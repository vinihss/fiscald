package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fiscald/internal/contingency"
	"fiscald/internal/domain"
	"fiscald/internal/repository"

	"github.com/google/uuid"
	"github.com/kardianos/service"
)

type program struct {
	server *http.Server
	repo   *repository.Repository
	cb     *contingency.CircuitBreaker
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	var err error
	p.repo, err = repository.NewSQLiteRepository("fiscald.db")
	if err != nil {
		log.Fatalf("Erro no banco de dados SQLite: %v", err)
	}

	p.cb = contingency.NewCircuitBreaker(3, 3*time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/emissao", p.handleEmissao)

	p.server = &http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	go p.workerRetransmissao()

	log.Println("[fiscald] Daemon iniciado na porta :9090")
	if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Erro no servidor HTTP: %v", err)
	}
}

func (p *program) Stop(s service.Service) error {
	log.Println("[fiscald] Encerrando o Daemon fiscald...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if p.server != nil {
		p.server.Shutdown(ctx)
	}
	if p.repo != nil {
		p.repo.Close()
	}
	return nil
}

func (p *program) handleEmissao(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo nao permitido", http.StatusMethodNotAllowed)
		return
	}

	var req domain.NFCeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}

	tpEmis := 1
	status := "AUTORIZADO"

	if p.cb.ShouldUseContingency() {
		tpEmis = 9
		status = "CONTINGENCIA"
	}

	chaveAcesso := fmt.Sprintf("3526091234567800019965001000000101%d1000001019", time.Now().UnixNano()%10)
	idRecord := uuid.New().String()

	rec := &repository.Record{
		ID:          idRecord,
		TenantID:    req.TenantID,
		VendaID:     req.VendaID,
		ChaveAcesso: chaveAcesso,
		NumeroNota:  101,
		Serie:       1,
		Status:      status,
		TpEmis:      tpEmis,
		XMLEnvio:    []byte("<NFe>...mock xml...</NFe>"),
	}

	if err := p.repo.Save(rec); err != nil {
		http.Error(w, "Erro ao salvar nota localmente", http.StatusInternalServerError)
		return
	}

	resp := domain.NFCeResponse{
		Status:      status,
		TipoEmissao: tpEmis,
		ChaveAcesso: chaveAcesso,
		QRCodeURL:   "https://www.sefaz.sp.gov.br/nfce/qrcode?p=" + chaveAcesso,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (p *program) workerRetransmissao() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		records, err := p.repo.GetPendingContingency()
		if err != nil || len(records) == 0 {
			continue
		}

		log.Printf("[fiscald-sync] Encontradas %d notas pendentes em contingencia. Processando...", len(records))
		for _, rec := range records {
			p.repo.UpdateStatus(rec.ID, "AUTORIZADO", []byte("<protNFe>...ok...</protNFe>"), "")
		}
	}
}

func main() {
	svcFlag := flag.String("service", "", "Acao do servico: install, uninstall, start, stop")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "fiscald",
		DisplayName: "Daemon de Integracao Fiscal NFC-e (Go)",
		Description: "Motor fiscal autonomo para emissao e contingencia de NFC-e em restaurantes.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatalf("Erro ao configurar servico: %v", err)
	}

	if *svcFlag != "" {
		err = service.Control(s, *svcFlag)
		if err != nil {
			log.Fatalf("Falha na operacao do servico (%s): %v", *svcFlag, err)
		}
		fmt.Printf("Operacao de servico '%s' executada com sucesso.\n", *svcFlag)
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		s.Stop()
		os.Exit(0)
	}()

	if err := s.Run(); err != nil {
		log.Fatalf("Erro ao executar daemon: %v", err)
	}
}
