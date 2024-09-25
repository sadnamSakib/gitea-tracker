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

func Repo(org string, repo model.Repo, users []model.User, dailyCommitCountListForAllUsers [][]int, dateListForAllUsers []string, viewBy string) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<body>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Navbar().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!-- Include Chart.js from CDN --><script src=\"https://cdn.jsdelivr.net/npm/chart.js\"></script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = RepoScript(org, repo, users, dailyCommitCountListForAllUsers, dateListForAllUsers, viewBy).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"container mx-auto py-10 px-10\"><div class=\"grid grid-cols-1 md:grid-cols-3 gap-6\"><div class=\"col-span-1 bg-white shadow-md rounded-lg px-4 py-2\"><div class=\"text-3xl text-gray-700 text-center mb-2\"><strong>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(repo.Name)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/repo.templ`, Line: 28, Col: 46}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</strong></div><div class=\"text-sm text-gray-700 text-center mb-2\"><strong>Created On :</strong> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(repo.Created.Format("02 January 2006"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/repo.templ`, Line: 32, Col: 97}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"overflow-y-auto rounded-lg\" style=\"max-height: 450px;\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for i, user := range users {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bg-white shadow-md rounded-lg p-4 m-1 hover:shadow-lg transition-shadow cursor-pointer flex flex-wrap items-center\" data-user-id=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(user.Username)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/repo.templ`, Line: 37, Col: 188}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><!-- User Avatar --><div class=\"flex-shrink-0 mr-4\"><img src=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(user.Avatar_url)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/repo.templ`, Line: 41, Col: 66}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" alt=\"Avatar of { user.Username }\" class=\"rounded-full w-12 h-12 object-cover\"></div><!-- Username and Email --><div><div class=\"text-lg font-semibold\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 string
			templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(user.Username)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/repo.templ`, Line: 46, Col: 90}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><!-- User Email --><div class=\"text-sm \">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var7 string
			templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(user.Email)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/repo.templ`, Line: 48, Col: 74}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div><!-- Commits Count (Placeholder) --><div class=\"ml-auto text-right\"><div class=\"text-sm font-semibold\">Commits: ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var8 string
			templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(func() string {
				total := 0
				for _, commitCount := range dailyCommitCountListForAllUsers[i] {
					total += commitCount
				}
				return strconv.Itoa(total)
			}())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/repo.templ`, Line: 59, Col: 47}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if len(users) == 0 {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"col-span-full p-4 text-red-500 text-center\">No users found.</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div><div class=\"col-span-2\"><div class=\"flex justify-end mb-2\"><button type=\"button\" class=\"bg-blue-500 hover:bg-blue-700 text-white font-bold py-1 px-2 rounded inline-flex items-center\"><svg class=\"w-4 h-4 mr-2\" xmlns=\"http://www.w3.org/2000/svg\" height=\"14\" width=\"12.25\" viewBox=\"0 0 448 512\"><!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.--><path fill=\"#FFD43B\" d=\"M9.4 233.4c-12.5 12.5-12.5 32.8 0 45.3l160 160c12.5 12.5 32.8 12.5 45.3 0s12.5-32.8 0-45.3L109.2 288 416 288c17.7 0 32-14.3 32-32s-14.3-32-32-32l-306.7 0L214.6 118.6c12.5-12.5 12.5-32.8 0-45.3s-32.8-12.5-45.3 0l-160 160z\"></path></svg> <a href=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var9 templ.SafeURL = templ.SafeURL(fmt.Sprintf("/orgs/%s/repos", org))
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var9)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">Back To Search</a></button></div><div class=\"flex justify-end items-center mb-4\"><!-- Right section for View By dropdown --><div class=\"text-right\"><label for=\"viewBy\" class=\"mr-2\">View By:</label> <select id=\"viewBy\" name=\"viewBy\" class=\"border p-1\">")
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</select></div></div><div class=\"mb-2\" widht=\"250\" height=\"250\"><canvas id=\"repoCommitChart\" width=\"200\" height=\"200\"></canvas></div><div class=\"text-left\"><div class=\"text-xl font-semibold\">Selected User : <span id=\"selectedUser\">All</span></div><div class=\"text-xl font-semibold\">Commits : <span id=\"totalCommits\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var10 string
		templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(func() int {
			total := 0
			for _, commitCountList := range dailyCommitCountListForAllUsers {
				for _, commitCount := range commitCountList {
					total += commitCount
				}
			}
			return total
		}()))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/repo.templ`, Line: 140, Col: 40}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div></div><div class=\"my-6\" widht=\"250\" height=\"250\"><canvas id=\"commitDailyChart\" width=\"200\" height=\"200\"></canvas></div></div></div></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func RepoScript(org string, repo model.Repo, users []model.User, dailyCommitCountListForAllUsers [][]int, dateListForAllUsers []string, viewBy string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_RepoScript_b36e`,
		Function: `function __templ_RepoScript_b36e(org, repo, users, dailyCommitCountListForAllUsers, dateListForAllUsers, viewBy){var lineChart = null;
    var piChart = null;
    var selectedUserId =null;
    function renderCommitDailyChart(dates, commitCounts) {

        if (lineChart) {
            lineChart.destroy();
        }
        var ctx = document.getElementById('commitDailyChart').getContext('2d');
        lineChart = new Chart(ctx, {
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
                maintainAspectRatio: false, 
                aspectRatio: 1,
                scales: {
                    y: {
                        beginAtZero: true
                    }
                },
                plugins: {
                    legend: {
                        display: false 
                    }
                }
            }
        });
    }

    function renderRepoCommitChart(repoNames, repoCommitCounts) {
        if(piChart){
            piChart.destroy();
        }
        var ctx = document.getElementById('repoCommitChart').getContext('2d');
        pChart = new Chart(ctx, {
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
                maintainAspectRatio: false,
                aspectRatio: 1,
                layout: {
                    padding: {
                        left: 0,
                        right: 30
                    }
                },
                plugins: {
                    legend: {
                        position: 'left',
                        align: 'start',
                        labels: {
                            boxWidth: 20,
                            padding: 10
                        }
                    }
                }
            }
        });
        
    }

    window.onload = function () {
        console.log(dailyCommitCountListForAllUsers);
        console.log(dateListForAllUsers);

        var aggregatedCommitCount = [];
        for (var i = 0; i < dateListForAllUsers.length; i++) {
            var total = 0;
            for (var j = 0; j < dailyCommitCountListForAllUsers.length; j++) {
                total += dailyCommitCountListForAllUsers[j][i];
            }
            aggregatedCommitCount.push(total);
        }

        // Render the initial aggregated chart for all users
        renderCommitDailyChart(dateListForAllUsers, aggregatedCommitCount);

        var aggregatedCommitCountListForUsers = [];
        for (var i = 0; i < users.length; i++) {
            var total = 0;
            for (var j = 0; j < dailyCommitCountListForAllUsers[i].length; j++) {
                total += dailyCommitCountListForAllUsers[i][j];
            }
            aggregatedCommitCountListForUsers.push(total);
        }
        renderRepoCommitChart(users.map(user => user.username), aggregatedCommitCountListForUsers);

        document.querySelectorAll('[data-user-id]').forEach((el) => {
            el.addEventListener('click', () => {
                const clickedUserId = el.getAttribute('data-user-id');
                
                // If the same user is clicked again, revert to aggregated commits
                if (selectedUserId === clickedUserId) {
                    // Reset selected user
                    selectedUserId = null;
                    
                   
                    el.style.backgroundColor = 'white';
                    el.style.color = 'black';
                    
                    // Reset to the aggregated commit chart
                    document.getElementById('selectedUser').innerText = "All";
                    document.getElementById('totalCommits').innerText = aggregatedCommitCount.reduce((a, b) => a + b, 0);
                    renderCommitDailyChart(dateListForAllUsers, aggregatedCommitCount);
                } else {
                    // Highlight the selected card
                    document.querySelectorAll('[data-user-id]').forEach(card => {
                        card.style.backgroundColor = 'white';
                        card.style.color = 'black';
                    });
                    el.style.backgroundColor = 'blue';
                    el.style.color = 'white';
                    
                    selectedUserId = clickedUserId;

                    // Get the commit counts for the clicked user
                    const index = users.findIndex(user => user.username === clickedUserId);

                    if (index !== -1) {
                        console.log('Selected User:', users[index].Username);

                        // Get the commit counts for the clicked user
                        const userCommitCounts = dailyCommitCountListForAllUsers[index];

                        document.getElementById('selectedUser').innerText = clickedUserId;
                        document.getElementById('totalCommits').innerText = userCommitCounts.reduce((a, b) => a + b, 0);
                        renderCommitDailyChart(dateListForAllUsers, userCommitCounts);
                    } else {
                        console.error('User not found');
                    }
                }
            });
        });

        document.getElementById('viewBy').addEventListener('change', function () {
            var viewBy = document.getElementById('viewBy').value;
            var url = ` + "`" + `/orgs/${org}/repos/${repo.name}?viewBy=${viewBy}` + "`" + `;
            window.location.href = url;  
        });
    };
    
}`,
		Call:       templ.SafeScript(`__templ_RepoScript_b36e`, org, repo, users, dailyCommitCountListForAllUsers, dateListForAllUsers, viewBy),
		CallInline: templ.SafeScriptInline(`__templ_RepoScript_b36e`, org, repo, users, dailyCommitCountListForAllUsers, dateListForAllUsers, viewBy),
	}
}

var _ = templruntime.GeneratedTemplate
