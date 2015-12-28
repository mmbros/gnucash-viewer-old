package xml

import (
	"compress/gzip"
	"encoding/xml"
	"os"
)

/*
<gnc-v2>
	<gnc:book>
		<gnc:account>
		<gnc:transaction>
			<trn:slots>
				<slot>
      <trn:splits>
        <trn:split>
*/

// Gnc type
type Gnc struct {
	XMLName xml.Name `xml:"gnc-v2"`
	Books   []Book   `xml:"book"`
}

// Book type
type Book struct {
	XMLName         xml.Name      `xml:"book"`
	ID              string        `xml:"id"`
	AccountList     []Account     `xml:"account"`
	TransactionList []Transaction `xml:"transaction"`
}

// Account type
type Account struct {
	ID          string `xml:"id"`
	Type        string `xml:"type"`
	Name        string `xml:"name"`
	Description string `xml:"description"`
	ParentID    string `xml:"parent"`
	Currency    string `xml:"commodity>id"`
}

// Split type
type Split struct {
	ID              string `xml:"id"`
	ReconciledState string `xml:"reconciled-state"`
	ReconcileDate   string `xml:"reconcile-date>date"`
	Value           string `xml:"value"`
	Quantity        string `xml:"quantity"`
	AccountID       string `xml:"account"`
}

// Transaction type
type Transaction struct {
	ID          string  `xml:"id"`
	Currency    string  `xml:"currency>id"`
	DatePosted  string  `xml:"date-posted>date"`
	DateEntered string  `xml:"date-entered>date"`
	Description string  `xml:"description"`
	SplitList   []Split `xml:"splits>split"`
}

/*
<slot>
  <slot:key>notes</slot:key>
  <slot:value type="string"></slot:value>
<slot>

<slot>
  <slot:key>options</slot:key>
  <slot:value type="frame">
    <slot>
      <slot:key>Budgeting</slot:key>
      <slot:value type="frame"/>
    </slot>
  </slot:value>
</slot>

<slot>
  <slot:key>1</slot:key>
  <slot:value type="numeric">800/100</slot:value>
</slot>

<slot>
  <slot:key>date-posted</slot:key>
  <slot:value type="gdate">
    <gdate>2013-04-11</gdate>
  </slot:value>
</slot>

<slot>
  <slot:key>import-map-bayes</slot:key>
  <slot:value type="frame">
    <slot>
      <slot:key>Arancio</slot:key>
      <slot:value type="frame">
        <slot>
          <slot:key>Attività:Attività correnti:Conto corrente</slot:key>
          <slot:value type="integer">1</slot:value>
        </slot>
      </slot:value>
    </slot>
  </slot:value>
</slot>

*/

// Slot type : integer | string | frame | gdate | numeric
type Slot struct {
	XMLName  xml.Name `xml:"slot"`
	Key      string   `xml:"key"`
	Value    string   `xml:"value"`
	Children []*Slot
}

// ReadFile read the gnucash file in XML format
func ReadFile(path string) (*Gnc, error) {

	// open gnucash file
	gnucashFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer gnucashFile.Close()

	// decompress gnucash file
	reader, err := gzip.NewReader(gnucashFile)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// unmarshall XML
	gnc := Gnc{}
	err = xml.NewDecoder(reader).Decode(&gnc)
	if err != nil {
		return nil, err
	}

	return &gnc, nil
}
