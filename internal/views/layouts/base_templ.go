// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package layouts

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Base(content templ.Component) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Domains Manager</title><script src=\"https://unpkg.com/htmx.org@1.9.10\"></script><script src=\"https://cdn.tailwindcss.com\"></script><script defer src=\"https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js\"></script><link rel=\"stylesheet\" href=\"https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css\"></head><body class=\"bg-gray-50\"><div class=\"min-h-screen\"><!-- Боковая навигация --><div class=\"fixed inset-y-0 left-0 z-50 w-64 bg-white border-r border-gray-200\"><div class=\"flex h-16 shrink-0 items-center px-6 border-b border-gray-200\"><span class=\"text-xl font-bold text-gray-900\">DM</span></div><nav class=\"flex flex-1 flex-col\"><ul role=\"list\" class=\"flex flex-1 flex-col gap-y-7 px-6 py-4\"><li><ul role=\"list\" class=\"-mx-2 space-y-1\"><li><a href=\"/domains\" class=\"group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6 text-gray-900 hover:bg-gray-50\"><i class=\"fa-solid fa-globe text-indigo-600 text-lg\"></i> <span>Domains</span></a></li><li><a href=\"/mailboxes\" class=\"group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6 text-gray-600 hover:bg-gray-50 hover:text-gray-900\"><i class=\"fa-solid fa-envelope text-indigo-600 text-lg\"></i> <span>Mailboxes</span></a></li><li><a href=\"/statistics\" class=\"group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6 text-gray-600 hover:bg-gray-50 hover:text-gray-900\"><i class=\"fa-solid fa-chart-line text-indigo-600 text-lg\"></i> <span>Statistics</span></a></li><li><a href=\"/logs\" class=\"group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6 text-gray-600 hover:bg-gray-50 hover:text-gray-900\"><i class=\"fa-solid fa-list-check text-indigo-600 text-lg\"></i> <span>Logs</span></a></li></ul></li><li class=\"mt-auto border-t border-gray-200 pt-4\"><a href=\"/settings\" class=\"group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6 text-gray-600 hover:bg-gray-50 hover:text-gray-900\"><i class=\"fa-solid fa-gear text-indigo-600 text-lg\"></i> <span>Settings</span></a></li></ul></nav></div><!-- Основной контент --><div class=\"pl-64\"><!-- Верхняя панель --><div class=\"sticky top-0 z-40 flex h-16 shrink-0 items-center gap-x-4 border-b border-gray-200 bg-white px-4 shadow-sm sm:gap-x-6 sm:px-6 lg:px-8\"><div class=\"flex flex-1 gap-x-4 self-stretch lg:gap-x-6\"><div class=\"flex flex-1\"></div><div class=\"flex items-center gap-x-4 lg:gap-x-6\"><button type=\"button\" class=\"-m-1.5 flex items-center p-1.5\"><span class=\"sr-only\">Open user menu</span> <i class=\"fa-solid fa-circle-user text-gray-400 hover:text-indigo-600 text-2xl transition-colors\"></i></button></div></div></div><!-- Основное содержимое --><main class=\"py-6 px-4 sm:px-6 lg:px-8\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = content.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</main></div></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
