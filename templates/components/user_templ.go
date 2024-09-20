// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"gitea.vivasoftltd.com/Vivasoft/gitea-commiter-plugin/pkg/model"
	"strconv"
)

var onceHandle = templ.NewOnceHandle()

func User(user model.User, dailyCommitCountList []int, repoCommitCountList []int, repoNameList []string, dateList []string, viewBy string) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Header().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<body><!-- Include Chart.js from CDN --><script src=\"https://cdn.jsdelivr.net/npm/chart.js\"></script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = UserScript(user.Username, dailyCommitCountList, repoCommitCountList, repoNameList, dateList).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"container mx-auto py-10\" style=\"min-height:100vh;\"><div class=\"grid grid-cols-1 md:grid-cols-3 gap-6\"><div class=\"col-span-1 bg-white shadow-md rounded-lg p-4\"><div class=\"flex justify-center mb-4\"><img src=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(user.Avatar_url)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/user.templ`, Line: 28, Col: 53}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" alt=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("Avatar of %s", user.Username))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/user.templ`, Line: 28, Col: 102}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"rounded-full w-32 h-32 object-cover\"></div><div class=\"text-lg font-semibold text-center mb-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(user.Username)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/user.templ`, Line: 32, Col: 90}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"text-sm text-gray-700 text-center mb-2\"><strong>Email:</strong> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(user.Email)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/user.templ`, Line: 36, Col: 63}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"text-sm text-gray-700 text-center mb-2\"><strong>Total Commits:</strong> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 string
		templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(user.Total_commits))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/user.templ`, Line: 40, Col: 93}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"text-sm text-gray-700 text-center mb-2\"><strong>Repositories:</strong> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var7 string
		templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(len(repoNameList)))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/user.templ`, Line: 45, Col: 91}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"mb-2\"><canvas id=\"repoCommitChart\" width=\"300\" height=\"300\"></canvas></div></div><div class=\"col-span-2\"><!-- Date Selection and Repo Selection Dropdown --><div class=\"flex justify-end mx-auto\"><div><label for=\"viewBy\">View By:</label> <select id=\"viewBy\" name=\"viewBy\" class=\"border p-1\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		switch viewBy {
		case "week":
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"week\" selected>This Week</option> <option value=\"month\">This Month</option> <option value=\"year\">This Year</option> <option value=\"allTime\">All Time</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case "month":
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"week\">This Week</option> <option value=\"month\" selected>This Month</option> <option value=\"year\">This Year</option> <option value=\"allTime\">All Time</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case "year":
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"week\">This Week</option> <option value=\"month\">This Month</option> <option value=\"year\" selected>This Year</option> <option value=\"allTime\">All Time</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case "allTime":
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"week\">This Week</option> <option value=\"month\">This Month</option> <option value=\"year\">This Year</option> <option value=\"allTime\" selected>All Time</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		default:
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"week\" selected>This Week</option> <option value=\"month\">This Month</option> <option value=\"year\">This Year</option> <option value=\"allTime\">All Time</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</select></div></div><div class=\"my-6\"><canvas id=\"commitDailyChart\" width=\"300\" height=\"300\"></canvas></div></div></div></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func UserScript(username string, dailyCommitCountList []int, repoCommitCountList []int, repoNameList []string, dateList []string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_UserScript_3f5b`,
		Function: `function __templ_UserScript_3f5b(username, dailyCommitCountList, repoCommitCountList, repoNameList, dateList){function renderCommitDailyChart(dates, commitCounts) {
        var ctx = document.getElementById('commitDailyChart').getContext('2d');
        new Chart(ctx, {
            type: 'line',
            data: {
                labels: dates,
                datasets: [{
                    label: 'Commits',
                    data: commitCounts,
                    fill: true,
                    borderColor: 'rgba(75, 192, 192, 1)',
                    tension: 0.5
                }]
            },
            options: {
                responsive: true, 
                maintainAspectRatio: false, // Allows you to control the aspect ratio
                scales: {
                    y: {
                        beginAtZero: true
                    }
                }
            }
        });
    }

    function renderRepoCommitChart(repoNames, repoCommitCounts) {
        var ctx = document.getElementById('repoCommitChart').getContext('2d');
        new Chart(ctx, {
            type: 'doughnut',
            data: {
                labels: repoNames,
                datasets: [{
                    label: 'Commits per Repo',
                    data: repoCommitCounts,
                    backgroundColor: [
                        'rgba(255, 99, 132, 0.2)',
                        'rgba(54, 162, 235, 0.2)',
                        'rgba(255, 206, 86, 0.2)',
                        'rgba(75, 192, 192, 0.2)',
                        'rgba(153, 102, 255, 0.2)',
                        'rgba(255, 159, 64, 0.2)'
                    ],
                    borderColor: [
                        'rgba(255, 99, 132, 1)',
                        'rgba(54, 162, 235, 1)',
                        'rgba(255, 206, 86, 1)',
                        'rgba(75, 192, 192, 1)',
                        'rgba(153, 102, 255, 1)',
                        'rgba(255, 159, 64, 1)'
                    ],
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false 
            }
        });
    }

    window.onload = function () {
        
        renderCommitDailyChart(dateList, dailyCommitCountList);
        renderRepoCommitChart(repoNameList, repoCommitCountList);

       
        document.getElementById('viewBy').addEventListener('change', function () {
            var viewBy = document.getElementById('viewBy').value;
                var url = ` + "`" + `/web/users/${username}?viewBy=${viewBy}` + "`" + `;
                window.location.href = url;
                
        });
    };
    
}`,
		Call:       templ.SafeScript(`__templ_UserScript_3f5b`, username, dailyCommitCountList, repoCommitCountList, repoNameList, dateList),
		CallInline: templ.SafeScriptInline(`__templ_UserScript_3f5b`, username, dailyCommitCountList, repoCommitCountList, repoNameList, dateList),
	}
}

var _ = templruntime.GeneratedTemplate
