package invoice_handlers

import (
	"net/http"
	"strings"

	"go-upcycle_connect-backend/app/middleware/auth_middleware"
	"go-upcycle_connect-backend/app/models/invoice_models"
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/utils/pdf"
	"go-upcycle_connect-backend/utils/request"
	"go-upcycle_connect-backend/utils/response"
)

// IndexInvoicesHandler — GET /invoices/me (auth)
// Liste les factures de l'utilisateur (achats payes + abonnements).
func IndexInvoicesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())
	response.NewSuccessData(w, invoice_models.GetUserInvoices(userId))
}

// InvoicePdfHandler — GET /invoices/{ref}/pdf (auth)
// Genere et renvoie la facture au format PDF (le "double conserve" = le
// paiement/abonnement sous-jacent, le PDF est genere a la demande).
func InvoicePdfHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
	userId := auth_middleware.GetUserId(r.Context())

	ref := request.Request(r, "ref").Value()
	if ref == "" {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusNotFound)
		return
	}

	inv := invoice_models.GetUserInvoiceByRef(userId, ref)
	if inv == nil {
		response.NewErrorMessage(w, response.ErrInvoiceNotFound, http.StatusNotFound)
		return
	}

	bytes, err := pdf.GenerateInvoice(*inv, userId)
	if err != nil {
		log.Info("generate invoice pdf: " + err.Error())
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}

	filename := strings.ToUpper(strings.ReplaceAll(ref, "_", "-")) + ".pdf"
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
}
