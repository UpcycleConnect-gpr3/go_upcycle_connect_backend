package pdf

import (
	"bytes"
	"fmt"
	"strings"

	"go-upcycle_connect-backend/app/models/invoice_models"

	"github.com/go-pdf/fpdf"
)

func euros(cents int) string {
	return fmt.Sprintf("%.2f EUR", float64(cents)/100)
}

func GenerateInvoice(inv invoice_models.Invoice, clientRef string) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 22)
	pdf.SetTextColor(30, 30, 30)
	pdf.Cell(0, 12, "UpcycleConnect")
	pdf.Ln(11)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(90, 90, 90)
	pdf.Cell(0, 5, "174, rue La Fayette - 75010 Paris")
	pdf.Ln(5)
	pdf.Cell(0, 5, "contact@upcycleconnect.fr")
	pdf.Ln(14)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(30, 30, 30)
	pdf.Cell(0, 10, "FACTURE")
	pdf.Ln(11)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(60, 60, 60)
	ref := strings.ToUpper(strings.ReplaceAll(inv.Ref, "_", "-"))
	pdf.Cell(0, 6, "Reference : "+ref)
	pdf.Ln(6)
	pdf.Cell(0, 6, "Date : "+dateOnly(inv.IssuedAt))
	pdf.Ln(6)
	pdf.Cell(0, 6, "Client : "+clientRef)
	pdf.Ln(6)
	pdf.Cell(0, 6, "Statut : "+statusLabel(inv.Status))
	pdf.Ln(14)

	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.SetTextColor(30, 30, 30)
	pdf.CellFormat(120, 9, "Designation", "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 9, "Montant", "1", 0, "R", true, 0, "")
	pdf.Ln(9)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(120, 9, tr(inv.Label), "1", 0, "L", false, 0, "")
	pdf.CellFormat(50, 9, euros(inv.AmountCents), "1", 0, "R", false, 0, "")
	pdf.Ln(9)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(120, 10, "Total TTC", "1", 0, "R", false, 0, "")
	pdf.CellFormat(50, 10, euros(inv.AmountCents), "1", 0, "R", false, 0, "")
	pdf.Ln(20)

	pdf.SetFont("Arial", "I", 9)
	pdf.SetTextColor(120, 120, 120)
	pdf.MultiCell(0, 5, tr("Document genere automatiquement par UpcycleConnect. TVA non applicable, art. 293 B du CGI (selon regime)."), "", "L", false)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func dateOnly(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}

func statusLabel(s string) string {
	switch s {
	case "paid":
		return "Paye"
	case "active":
		return "Actif"
	default:
		return s
	}
}

func tr(s string) string {
	r := strings.NewReplacer(
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"à", "a", "â", "a", "î", "i", "ï", "i",
		"ô", "o", "ö", "o", "û", "u", "ü", "u", "ç", "c",
		"É", "E", "È", "E", "À", "A", "—", "-", "’", "'",
		"CO₂", "CO2",
	)
	return r.Replace(s)
}
