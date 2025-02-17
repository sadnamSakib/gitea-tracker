// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Navbar() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<nav class=\"bg-white sticky top-0 shadow-md \"><div class=\"container mx-auto px-4\"><div class=\"flex justify-between items-center py-4\"><div class=\"text-2xl font-bold text-gray-800 flex items-center justify-center\"><!-- Entire logo and text are clickable --><a href=\"/\" class=\"flex items-between\"><!-- Logo on the left --><img src=\"/web/assets/vivalogo.png\" alt=\"Logo\" class=\"mr-3\" style=\"height: 1em; width: 2em;\"> Gitea Tracker</a></div><div class=\"hidden md:flex space-x-6\"><a href=\"/orgs\" class=\"text-gray-700 hover:text-gray-900 font-semibold\">Repositories</a> <a href=\"/users\" class=\"text-gray-700 hover:text-gray-900 font-semibold\">Users</a></div><!-- Mobile Menu Button --><div class=\"md:hidden\"><button id=\"menu-toggle\" class=\"text-gray-700 focus:outline-none\"><svg class=\"w-6 h-6\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\" xmlns=\"http://www.w3.org/2000/svg\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M4 6h16M4 12h16m-7 6h7\"></path></svg></button></div></div></div><div id=\"mobile-menu\" class=\"md:hidden bg-white px-4 pt-2 pb-4 space-y-2 hidden\"><a href=\"/orgs\" class=\"block text-gray-700 hover:text-gray-900 font-semibold\">Repositories</a> <a href=\"/users\" class=\"block text-gray-700 hover:text-gray-900 font-semibold\">Users</a></div></nav>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func NavbarScript() templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_NavbarScript_c017`,
		Function: `function __templ_NavbarScript_c017(){// Toggle mobile menu visibility
        document.getElementById('menu-toggle').addEventListener('click', function () {
          var mobileMenu = document.getElementById('mobile-menu');
          mobileMenu.classList.toggle('hidden');
        });
    
}`,
		Call:       templ.SafeScript(`__templ_NavbarScript_c017`),
		CallInline: templ.SafeScriptInline(`__templ_NavbarScript_c017`),
	}
}

var _ = templruntime.GeneratedTemplate
