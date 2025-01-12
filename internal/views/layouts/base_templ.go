// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package layouts

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/danilshap/domains-generator/internal/views/components/common"

func Base(contents templ.Component) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html class=\"h-full bg-gray-100\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Domains Manager</title><script src=\"https://unpkg.com/htmx.org@1.9.6\"></script><script src=\"https://cdn.tailwindcss.com\"></script><script defer src=\"https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js\"></script><link rel=\"stylesheet\" href=\"https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css\"><link href=\"https://cdn.jsdelivr.net/npm/tom-select@2.3.1/dist/css/tom-select.css\" rel=\"stylesheet\"><script src=\"https://cdn.jsdelivr.net/npm/tom-select@2.3.1/dist/js/tom-select.complete.min.js\"></script><script src=\"/assets/js/websocket.js\"></script><script src=\"/assets/js/notifications.js\"></script></head><body class=\"h-full\"><div class=\"min-h-full\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
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
			templ_7745c5c3_Err = contents.Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = common.Nav().Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = common.Modal().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><script>\n\t\t\t\tfunction validateDomainForm(form) {\n\t\t\t\t\tconst domain = form.name.value;\n\t\t\t\t\tconst pattern = /^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\\.[a-zA-Z]{2,}$/;\n\t\t\t\t\tif (!pattern.test(domain)) {\n\t\t\t\t\t\talert(\"Please enter a valid domain name\");\n\t\t\t\t\t\treturn false;\n\t\t\t\t\t}\n\t\t\t\t\treturn true;\n\t\t\t\t}\n\n\t\t\t\tfunction updateEmailAddress() {\n\t\t\t\t\tconst localPart = document.getElementById('local_part').value;\n\t\t\t\t\tconst domainSelect = document.getElementById('domain_id');\n\t\t\t\t\tconst selectedOption = domainSelect.options[domainSelect.selectedIndex];\n\t\t\t\t\tconst domain = selectedOption.getAttribute('data-domain');\n\t\t\t\t\t\n\t\t\t\t\tif (localPart && domain) {\n\t\t\t\t\t\tdocument.getElementById('address').value = `${localPart}@${domain}`;\n\t\t\t\t\t}\n\t\t\t\t}\n\n\t\t\t\tfunction showModal() {\n\t\t\t\t\tdocument.getElementById('modal').style.display = 'block';\n\t\t\t\t}\n\n\t\t\t\tfunction closeModal() {\n\t\t\t\t\tdocument.getElementById('modal').style.display = 'none';\n\t\t\t\t\tdocument.querySelector('#modal > div:last-child > div > div').innerHTML = '';\n\t\t\t\t}\n\n\t\t\t\tdocument.body.addEventListener('htmx:afterOnLoad', function(evt) {\n\t\t\t\t\tif (evt.detail.target.id === 'modal') {\n\t\t\t\t\t\tshowModal();\n\t\t\t\t\t}\n\t\t\t\t});\n\n\t\t\t\tdocument.body.addEventListener('htmx:trigger', function(evt) {\n\t\t\t\t\tif (evt.detail.trigger === 'refreshMailboxes') {\n\t\t\t\t\t\thtmx.ajax('GET', '/mailboxes', '#mailboxes-list');\n\t\t\t\t\t}\n\t\t\t\t});\n\n\t\t\t\tdocument.addEventListener('htmx:afterSettle', function(evt) {\n\t\t\t\t\tevt.detail.elt.querySelectorAll('.tomselect').forEach(function(el){\n\t\t\t\t\t\tnew TomSelect(el, {\n\t\t\t\t\t\t\tcreate: false,\n\t\t\t\t\t\t\tsortField: {\n\t\t\t\t\t\t\t\tfield: \"text\",\n\t\t\t\t\t\t\t\tdirection: \"asc\"\n\t\t\t\t\t\t\t}\n\t\t\t\t\t\t});\n\t\t\t\t\t});\n\t\t\t\t});\n\t\t\t</script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
