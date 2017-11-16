package main

type Client struct {
	ClientId      string   `json:"clientId"`
	ClientName    string   `json:"clientName"`
	Policies      []string `json:"policies"`
	RFQArray      []string `json:"rfqArray"`
	ProposalArray []string `json:"proposalArray"`
}

type Insurer struct {
	InsurerId     string   `json:"insurerId"`
	InsurerName   string   `json:"insurerName"`
	RFQArray      []string `json:"rfqArray"`
	Quotes        []string `json:"quotes"`
	Policies      []string `json:"policies"`
	ProposalArray []string `json:"proposalArray"`
}

type RFQ struct {
	RFQId              string              `json:"rfqId"`
	ClientId           string              `json:"insurerId"`
	InsuredName        string              `json:"insuredName"`
	TypeOfInsurance    string              `json:"typeOfInsurance"`
	RiskAmount         float64             `json:"riskAmount"`
	RiskLocation       string              `json:"riskLocation"`
	StartDate          string              `json:"startDate"`
	EndDate            string              `json:"endDate"`
	Status             string              `json:"status"`
	Quotes             []string            `json:"quotes"`
	LeadQuote          string              `json:"leadQuote"`
	SelectedInsurer    []string            `json:"selectedInsurer"`
	LeadInsurer        string              `json:"leadInsurer"`
	ProposalDocHash    string              `json:"proposalDocHash"`
	ProposalNum        string              `json:"proposalNum"`
	Intermediary       string              `json:"intermediary"`
	TransactionHistory []TransactionRecord `json:"transactionHistory"`
}

type TransactionRecord struct {
	TxId      string `json:"txId"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type Proposal struct {
	ProposalNum        string              `json:"proposalNum"`
	RFQId              string              `json:"rfqId"`
	Status             string              `json:"status"`
	PolicyNum          string              `json:"policyNum"`
	TransactionHistory []TransactionRecord `json:"transactionHistory"`
}

type Policy struct {
	PolicyNumber       string              `json:"policyNumber"`
	ProposalNum        string              `json:"proposalNum"`
	PolicyDocHash      string              `json:"policyDocHash"`
	Details            RFQ                 `json:"details"`
	Status             string              `json:"status"`
	TransactionHistory []TransactionRecord `json:"transactionHistory"`
}

type Broker struct {
	BrokerId      string   `json:"brokerId"`
	BrokerName    string   `json:"brokerName"`
	Clients       []string `json:"clients"`
	RFQArray      []string `json:"rfqArray"`
	Policies      []string `json:"policies"`
	ProposalArray []string `json:"proposalArray"`
}

type Quote struct {
	QuoteId            string              `json:"quoteId"`
	InsurerName        string              `json:"insurerName"`
	InsurerId          string              `json:"insurerId"`
	Premium            float64             `json:"premium"`
	Capacity           float64             `json:"capacity"`
	RFQId              string              `json:"rfqId"`
	Status             string              `json:"status"`
	TransactionHistory []TransactionRecord `json:"transactionHistory"`
}
