// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package domains

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Form() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity\" id=\"modal-backdrop\"></div><div class=\"fixed inset-0 z-10 overflow-y-auto\"><div class=\"flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0\"><div class=\"relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6\"><div class=\"absolute right-0 top-0 hidden pr-4 pt-4 sm:block\"><button type=\"button\" class=\"rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none\" onclick=\"document.getElementById(&#39;modal-backdrop&#39;).remove(); this.closest(&#39;#domain-form&#39;).innerHTML = &#39;&#39;\"><span class=\"sr-only\">Close</span> <i class=\"fas fa-times h-6 w-6\"></i></button></div><div class=\"sm:flex sm:items-start\"><div class=\"mt-3 text-center sm:ml-4 sm:mt-0 sm:text-left w-full\"><h3 class=\"text-base font-semibold leading-6 text-gray-900\">Add New Domain</h3><div class=\"mt-4\"><form hx-post=\"/domains\" hx-target=\"#domains-list\" hx-swap=\"innerHTML\" hx-on::after-request=\"document.getElementById(&#39;modal-backdrop&#39;).remove(); this.closest(&#39;#domain-form&#39;).innerHTML = &#39;&#39;\" class=\"space-y-4\"><div><label for=\"name\" class=\"block text-sm font-medium text-gray-700\">Domain Name</label> <input type=\"text\" name=\"name\" id=\"name\" required class=\"mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm\" placeholder=\"example.com\"></div><div><label for=\"provider\" class=\"block text-sm font-medium text-gray-700\">Provider</label> <input type=\"text\" name=\"provider\" id=\"provider\" required class=\"mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm\" placeholder=\"Provider name\"></div><div class=\"mt-5 sm:mt-4 sm:flex sm:flex-row-reverse\"><button type=\"submit\" class=\"inline-flex w-full justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 sm:ml-3 sm:w-auto\">Save</button> <button type=\"button\" class=\"mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto\" onclick=\"document.getElementById(&#39;modal-backdrop&#39;).remove(); this.closest(&#39;#domain-form&#39;).innerHTML = &#39;&#39;\">Cancel</button></div></form></div></div></div></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
