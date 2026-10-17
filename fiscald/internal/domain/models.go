package domain

import "encoding/xml"

type NFCeRequest struct {
	TenantID   string          `json:"tenant_id"`
	VendaID    string          `json:"venda_id"`
	CPFCliente string          `json:"cpf_cliente,omitempty"`
	Itens      []ItemRequest   `json:"itens"`
	Pagamentos []PagamentoReq  `json:"pagamentos"`
}

type ItemRequest struct {
	Item          int     `json:"item"`
	Codigo        string  `json:"codigo"`
	Descricao     string  `json:"descricao"`
	NCM           string  `json:"ncm"`
	CFOP          string  `json:"cfop"`
	Unidade       string  `json:"unidade"`
	Quantidade    float64 `json:"quantidade"`
	ValorUnitario float64 `json:"valor_unitario"`
	ValorTotal    float64 `json:"valor_total"`
	ICMSOrigem    string  `json:"icms_origem"`
	ICMSCSOSN     string  `json:"icms_csosn"`
}

type PagamentoReq struct {
	Forma string  `json:"forma"`
	Valor float64 `json:"valor"`
}

type NFCeResponse struct {
	Status          string `json:"status"`
	TipoEmissao     int    `json:"tipo_emissao"`
	ChaveAcesso     string `json:"chave_acesso"`
	Protocolo       string `json:"protocolo,omitempty"`
	DataAutorizacao string `json:"data_autorizacao,omitempty"`
	QRCodeURL       string `json:"qr_code_url"`
	XMLPath         string `json:"xml_path,omitempty"`
	Motivo          string `json:"motivo,omitempty"`
}

type NFe struct {
	XMLName xml.Name `xml:"NFe"`
	Xmlns   string   `xml:"xmlns,attr"`
	InfNFe  InfNFe   `xml:"infNFe"`
}

type InfNFe struct {
	ID     string `xml:"Id,attr"`
	Versao string `xml:"versao,attr"`
	Ide    Ide    `xml:"ide"`
	Emit   Emit   `xml:"emit"`
	Dest   *Dest  `xml:"dest,omitempty"`
	Det    []Det  `xml:"det"`
	Total  Total  `xml:"total"`
	Pag    Pag    `xml:"pag"`
	InfAdic *InfAdic `xml:"infAdic,omitempty"`
}

type Ide struct {
	CUF     string `xml:"cUF"`
	CNF     string `xml:"cNF"`
	NatOp   string `xml:"natOp"`
	Mod     string `xml:"mod"`
	Serie   string `xml:"serie"`
	NNF     string `xml:"nNF"`
	DhEmi   string `xml:"dhEmi"`
	TpNF    string `xml:"tpNF"`
	IdDest  string `xml:"idDest"`
	CMunFG  string `xml:"cMunFG"`
	TpImp   string `xml:"tpImp"`
	TpEmis  string `xml:"tpEmis"`
	CDV     string `xml:"cDV"`
	TpAmb   string `xml:"tpAmb"`
	FinNFe  string `xml:"finNFe"`
	IndFinal string `xml:"indFinal"`
	IndPres string `xml:"indPres"`
	ProcEmi string `xml:"procEmi"`
	VerProc string `xml:"verProc"`
}

type Emit struct {
	CNPJ  string `xml:"CNPJ"`
	XNome string `xml:"xNome"`
	XFant string `xml:"xFant,omitempty"`
	EnderEnderEnderEnder struct {
		XLgr    string `xml:"xLgr"`
		Nro     string `xml:"nro"`
		XBairro string `xml:"xBairro"`
		CMun    string `xml:"cMun"`
		XMun    string `xml:"xMun"`
		UF      string `xml:"UF"`
		CEP     string `xml:"CEP"`
	} `xml:"enderEmit"`
	IE   string `xml:"IE"`
	CRT  string `xml:"CRT"`
}

type Dest struct {
	CPF string `xml:"CPF,omitempty"`
	CNPJ string `xml:"CNPJ,omitempty"`
	IndIEDest string `xml:"indIEDest"`
}

type Det struct {
	NItem int     `xml:"nItem,attr"`
	Prod  Prod    `xml:"prod"`
	Imposto Imposto `xml:"imposto"`
}

type Prod struct {
	CProd    string  `xml:"cProd"`
	CEAN     string  `xml:"cEAN"`
	XProd    string  `xml:"xProd"`
	NCM      string  `xml:"NCM"`
	CFOP     string  `xml:"CFOP"`
	UCom     string  `xml:"uCom"`
	QCom     float64 `xml:"qCom"`
	VUnCom   float64 `xml:"vUnCom"`
	VProd    float64 `xml:"vProd"`
	CEANTrib string  `xml:"cEANTrib"`
	UTrib    string  `xml:"uTrib"`
	QTrib    float64 `xml:"qTrib"`
	VUnTrib  float64 `xml:"vUnTrib"`
	IndTot   string  `xml:"indTot"`
}

type Imposto struct {
	ICMS ICMS `xml:"ICMS"`
}

type ICMS struct {
	ICMSSN102 *ICMSSN102 `xml:"ICMSSN102,omitempty"`
}

type ICMSSN102 struct {
	Orig   string `xml:"orig"`
	CSOSN  string `xml:"CSOSN"`
}

type Total struct {
	ICMSTot ICMSTot `xml:"ICMSTot"`
}

type ICMSTot struct {
	VBC     float64 `xml:"vBC"`
	VICMS   float64 `xml:"vICMS"`
	VICMSDeson float64 `xml:"vICMSDeson"`
	VFCP    float64 `xml:"vFCP"`
	VBCST   float64 `xml:"vBCST"`
	VST     float64 `xml:"vST"`
	VFCPST  float64 `xml:"vFCPST"`
	VFCPSTRet float64 `xml:"vFCPSTRet"`
	VProd   float64 `xml:"vProd"`
	VFrete  float64 `xml:"vFrete"`
	VSeg    float64 `xml:"vSeg"`
	VDesc   float64 `xml:"vDesc"`
	VII     float64 `xml:"vII"`
	VIPI    float64 `xml:"vIPI"`
	VIPIDevol float64 `xml:"vIPIDevol"`
	VPIS    float64 `xml:"vPIS"`
	VCOFINS float64 `xml:"vCOFINS"`
	VOutro  float64 `xml:"vOutro"`
	VNF     float64 `xml:"vNF"`
	VTottrib float64 `xml:"vTotTrib"`
}

type Pag struct {
	DetPag []DetPag `xml:"detPag"`
}

type DetPag struct {
	TPag string  `xml:"tPag"`
	VPag float64 `xml:"vPag"`
}

type InfAdic struct {
	InfCpl string `xml:"infCpl,omitempty"`
}
