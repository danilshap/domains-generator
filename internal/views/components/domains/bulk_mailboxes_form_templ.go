// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package domains

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	db "github.com/danilshap/domains-generator/internal/db/sqlc"
)

func BulkMailboxesForm(domain db.GetDomainByIDRow) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity\" id=\"modal-backdrop\"></div><div class=\"fixed inset-0 z-10 overflow-y-auto\"><div class=\"flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0\"><div class=\"relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6\"><div class=\"absolute right-0 top-0 hidden pr-4 pt-4 sm:block\"><button type=\"button\" class=\"rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none\" onclick=\"closeModal()\"><i class=\"fa-solid fa-xmark text-xl\"></i></button></div><div class=\"mb-5\"><h3 class=\"text-lg font-medium text-gray-900\">Create Multiple Mailboxes</h3><p class=\"mt-1 text-sm text-gray-500\">Create up to 100 mailboxes with a common prefix</p></div><form class=\"space-y-6\" hx-post=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("/domains/%d/bulk-mailboxes", domain.ID))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/views/components/domains/bulk_mailboxes_form.templ`, Line: 28, Col: 67}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#modal\"><input type=\"hidden\" name=\"domain_id\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(domain.ID))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/views/components/domains/bulk_mailboxes_form.templ`, Line: 31, Col: 72}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div><label for=\"prefix\" class=\"block text-sm font-medium text-gray-700\">Email Prefix</label><div class=\"relative mt-1\"><input type=\"text\" name=\"prefix\" id=\"prefix\" required class=\"block w-full rounded-lg border border-gray-300 bg-white px-4 py-2.5 pr-10 text-sm shadow-sm transition-colors hover:border-gray-400 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500\" placeholder=\"user\"><div class=\"pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3\"><i class=\"fa-solid fa-at text-gray-400 text-sm\"></i></div></div><p class=\"mt-1 text-xs text-gray-500\">Mailboxes will be created as: prefix.random@")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(domain.Name)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/views/components/domains/bulk_mailboxes_form.templ`, Line: 47, Col: 101}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p></div><div><label for=\"password\" class=\"block text-sm font-medium text-gray-700\">Password</label><div class=\"relative mt-1\"><input type=\"password\" name=\"password\" id=\"password\" required class=\"block w-full rounded-lg border border-gray-300 bg-white px-4 py-2.5 pr-10 text-sm shadow-sm transition-colors hover:border-gray-400 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500\"><div class=\"pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3\"><i class=\"fa-solid fa-lock text-gray-400 text-sm\"></i></div></div><p class=\"mt-1 text-xs text-gray-500\">All mailboxes will use this password</p></div><div class=\"flex justify-end space-x-3\"><button type=\"button\" class=\"px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50\" onclick=\"closeModal()\">Cancel</button> <button type=\"submit\" class=\"px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-lg hover:bg-indigo-700\">Create Mailboxes</button></div></form></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
