// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.598
package partial

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "example/hello/web/templates/components"

func Login(err string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex w-screen justify-center items-center h-screen bg-gray-800 \"><div class=\"flex flex-col max-w-md p-6 rounded-md sm:p-10 bg-gray-50 text-gray-800\"><div class=\"mb-8 text-center\"><h1 class=\"my-3 text-4xl font-bold\">Sign in</h1><p class=\"text-sm text-gray-600\">Sign in to access your account</p></div><form action=\"\" hx-post=\"\" method=\"POST\" class=\"space-y-12\"><div class=\"space-y-4\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.SignupError(err).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div><label for=\"email\" class=\"block mb-2 text-sm\">Email address</label> <input type=\"email\" name=\"email\" id=\"email\" placeholder=\"leroy@jenkins.com\" class=\"w-full px-3 py-2 border rounded-md border-gray-300 bg-gray-50 text-gray-800\"></div><div><div class=\"flex justify-between mb-2\"><label for=\"password\" class=\"text-sm\">Password</label> <a rel=\"noopener noreferrer\" href=\"/forgot-password\" class=\"text-xs hover:underline text-gray-600\">Forgot password?</a></div><input type=\"password\" name=\"password\" id=\"password\" placeholder=\"*****\" class=\"w-full px-3 py-2 border rounded-md border-gray-300 bg-gray-50 text-gray-800\"></div></div><div class=\"space-y-2\"><div><input type=\"submit\" class=\"w-full px-8 py-3 font-semibold rounded-md bg-violet-600 text-gray-50\" value=\"Login\"></div><p class=\"px-6 text-sm text-center text-gray-600\">Don't have an account yet? <a rel=\"noopener noreferrer\" href=\"/signup\" class=\"hover:underline text-violet-600\">Sign up</a>.</p></div></form></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
