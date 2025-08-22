package front

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type PageData struct {
	Result string
}

type HandlersHolder struct {
	Handlers map[string]func(http.ResponseWriter, *http.Request)
}

func CreatePageHandlers() HandlersHolder {
	tmpl := template.Must(template.New("page").Parse(`
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8" />
    <title>Форма отправки</title>
</head>
<body>
    <h1>Отправить строку</h1>
    <form id="myForm">
        <input type="text" name="inputStr" id="inputStr" placeholder="Введите строку" required />
        <button type="submit">Отправить</button>
    </form>
    <h2>Результат:</h2>
    <pre id="result"></pre>

    <script>
        document.getElementById('myForm').addEventListener('submit', function(e) {
            e.preventDefault();

            const inputVal = document.getElementById('inputStr').value;

            const url = 'http://localhost:8081/order/' + encodeURIComponent(inputVal);

            fetch(url)
                .then(response => {
                    if (!response.ok) {
                        throw new Error('не найдено');
                    }
                    return response.json();
                })
                .then(data => {
                    document.getElementById('result').textContent = JSON.stringify(data, null, 2);
                })
                .catch(error => {
                    document.getElementById('result').textContent = error.message;
                });
        });
    </script>
</body>
</html>
`))

	getHandler := func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}

	postHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var reqData struct {
			InputStr string `json:"inputStr"`
		}

		err := json.NewDecoder(r.Body).Decode(&reqData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp := map[string]string{
			"received": reqData.InputStr,
			"status":   "успешно получено",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
	holder := HandlersHolder{
		Handlers: make(map[string]func(http.ResponseWriter, *http.Request)),
	}
	holder.Handlers["/"] = getHandler
	holder.Handlers["/submit"] = postHandler
	return holder
}
