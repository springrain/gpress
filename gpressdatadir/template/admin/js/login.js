
//convert FormData to Object
var serializeFormToObject = function (form) {
    var objForm = {};
    var formData = new FormData(form);
    for (var key of formData.keys()) {
        objForm[key] = formData.get(key);
   }
    return objForm;
};
(function() {
    var loginForm = document.querySelector("#loginForm");
    loginForm.addEventListener("submit", function (event) {
        event.preventDefault();
        fetch(loginForm.action, {
            method: loginForm.method,
            //body: JSON.stringify(serializeFormToObject(loginForm)),
            credentials:"same-origin",
            mode:"same-origin",
            body: new FormData(loginForm)
        }).then(function (response) {
                if (response.ok) {
                    return response.json();
                }
            return Promise.reject(response);
        }).then(function(data) {
            console.log(data);
        }).catch(function (error) {
            console.warn(error);
        });
    });
})();