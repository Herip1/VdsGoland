function Main() {
    this.__proto__ = new Util();
    const util = this.__proto__;

    this.init = () => {
        console.log("init");
        document.querySelectorAll(".action").forEach(el => {
            console.log(el);
            el.onclick = () => {
                this[el.dataset.action]();
            };
        });
        this.checkTokenValidity();
    };

    this.checkTokenValidity = () => {
        const token = localStorage.getItem("token");
        //console.log(token)
        // Отправляем запрос на сервер для проверки токена
        util.get("/validateToken?token=" + token, data => {
            console.log(data, "saflknadsfl")
            if (data == true) {
                // Токен валиден, отображаем таблицу пользователей
                this.showInterface();
            } else {
                // Токен невалиден, от ображаем форму авторизации
                this.showAuthForm();
            }
        });
    };

    this.showAuthForm = () => {
        util.id("authForm").style.display = "block";
        util.id("usersTable").style.display = "none";
        util.id("vdsTable").style.display ="none";
        util.id("addUserBtn").style.display = "none";
        util.id("deleteUserBtn").style.display = "none";
        util.id("logoutBtn").style.display = "none";
    };

/*    this.showUsersTable = () => {
        util.id("authForm").style.display = "none";
        console.log("дошло не дошло")
        util.id("usersTable").style.display = "block";
        util.id("addUserBtn").style.display = "inline-block";
        util.id("deleteUserBtn").style.display = "inline-block";
        util.id("logoutBtn").style.display = "block";
        util.id("vdsTable").style.display = "block";
        util.id("addVdsBtn").style.display = "inline-block";// Показываем кнопку добавления VDS
        this.getUsers(); // Получаем и отображаем пользователей
        this.getVds()
    };
 */
    this.showInterface = () => {
        util.id("authForm").style.display = "none";
        util.id("usersTable").style.display = "block";
        util.id("addUserBtn").style.display = "inline-block";
        util.id("deleteUserBtn").style.display = "inline-block";
        util.id("logoutBtn").style.display = "block";
        util.id("vdsTable").style.display = "block";
        util.id("addVdsBtn").style.display = "inline-block";
        util.id("deleteVdsBtn").style.display = "inline-block";
        this.getUsers(); // Получаем и отображаем пользователей
        this.getVds()
        document.getElementById('showPaymentInfoBtn').style.display = 'block';
        document.getElementById('showPaymentInfoBtn').addEventListener('click', function() {
            var paymentInfo = document.getElementById('paymentInfo');
            paymentInfo.style.display = (paymentInfo.style.display === 'none' ? 'block' : 'none');
        });
    };

    this.hideInterface = () => {
        util.id("authForm").style.display = "none";
        util.id("usersTable").style.display = "none";
        util.id("addUserBtn").style.display = "none";
        util.id("deleteUserBtn").style.display = "none";
        util.id("deleteVdsBtn").style.display = "none";
        util.id("logoutBtn").style.display = "block";
        util.id("vdsTable").style.display = "block";
        util.id("addVdsBtn").style.display = "none";
        this.getVds()
        document.getElementById('showPaymentInfoBtn').style.display = 'block';
        document.getElementById('showPaymentInfoBtn').addEventListener('click', function() {
            var paymentInfo = document.getElementById('paymentInfo');
            paymentInfo.style.display = (paymentInfo.style.display === 'none' ? 'block' : 'none');
        });
    };

    this.addUser = () => {

        util.modal('addUserModal', 'show');

        document.getElementById('saveUserBtn').onclick = () => {
            const user = {
                login: util.id('addUserLogin').value,
                pass: util.id('addUserPassword').value,
                name: util.id('addUserName').value,
                access: util.id('addUserAccess').value === 'true'
            };

            util.post('/addUser', user, (response) => {
                if (response.success) {
                    this.getUsers();
                } else {
                    alert(response.message);
                }
            });

            util.modal('addUserModal', 'hide');
        };
    };


    this.deleteUser = () => {
        // Вызываем модальное окно
        util.modal('confirmDeleteModal', 'show');

        // Получаем кнопку подтверждения удаления и добавляем обработчик событий
        document.getElementById('confirmDeleteBtn').onclick = () => {
            const userId = parseInt(document.getElementById('deleteUserIdInput').value, 10);

            util.post('/deleteUser', {uid: userId}, (response) => {

                // Обработка ответа сервера
                if (response.success) {
                    // Обновить таблицу пользователей
                    this.getUsers();
                } else {
                    // Ошибка при удалении пользователя
                    alert(response.message);
                }
            });

            // Закрыть модальное окно
            util.modal('confirmDeleteModal', 'hide');
        };
    };

    this.deleteVds = () => {
        // Вызываем модальное окно
        util.modal('confirmDeleteModalVds', 'show');

        // Получаем кнопку подтверждения удаления и добавляем обработчик событий
        document.getElementById('confirmDeleteVdsBtn').onclick = () => {
            const vdsId = parseInt(document.getElementById('deleteVdsIdInput').value, 10);

            util.post('/deleteVds', {vid: vdsId}, (response) => {

                // Обработка ответа сервера
                if (response.success) {
                    // Обновить таблицу пользователей
                    this.getVds();
                } else {
                    // Ошибка при удалении пользователя
                    alert(response.message);
                }
            });

            // Закрыть модальное окно
            util.modal('confirmDeleteModalVds', 'hide');
        };
    };

    this.getUsers = () => {

        util.get("/getUsersTable", data1 => {
            console.log(data1, "Data from /getUsersTable");

            if (Array.isArray(data1)) {
                const tableBody = util.id("usersTable").querySelector("tbody");
                tableBody.innerHTML = ""; // Очищаем текущие строки

                data1.forEach(user => {
                    const row = `<tr>
                    <td>${user.uid}</td>
                    <td>${user.Login}</td>
                    <td>${user.Pass}</td>
                    <td>${user.Name}</td>
                    <td>${user.Access ? 'Да' : 'Нет'}</td>  
                </tr>`;
                    tableBody.innerHTML += row;
                });
            } else {
                console.error("Data is not an array:", data1);
            }
        }, error => {
            console.error("Error fetching data:", error);
        });
    };

    this.addVds = () => {
        util.modal('addVdsModal', 'show');

        document.getElementById('saveVdsBtn').onclick = () => {
            const vds = {
                name: util.id('addVdsName').value,
                img: util.id('addVdsImg').value,
                ram: parseInt(util.id('addVdsRam').value),
                hdd: parseInt(util.id('addVdsHdd').value),
                core: parseInt(util.id('addVdsCore').value),

            };

            util.post('/addVds', vds, (response) => {
                if (response.success) {
                    console.log("Vds added succefuly")
                    this.getVds();
                } else {
                    console.error('Error adding VDS', response.message)
                    alert(response.message);
                }
            });

            util.modal('addVdsModal', 'hide');
        };
        };

    this.getVds = () => {
        util.get("/getVdsTable", data2 => {
            console.log(data2, "Data from /getVdsTable");

            if (Array.isArray(data2)) {
                const tableBody = util.id("vdsTable").querySelector("tbody");
                tableBody.innerHTML = ""; // Очищаем текущие строки

                data2.forEach(vds => {
                    const row = `<tr>
                <td>${vds.uid}</td>
                <td>${vds.Name}</td>
                <td>${vds.Img}</td>
                <td>${vds.Ram}</td>
                <td>${vds.Hdd}</td>
                <td>${vds.Core}</td>
            </tr>`;
                    tableBody.innerHTML += row;
                });
            } else {
                console.error("Data is not an array:", data2);
            }
        }, error => {
            console.error("Error fetching data:", error);
        });
    };


    this.authIn = () => {
        //console.log("flskndlkfdsn11")
        util.post("/auth", {
            login: util.id("authLogin").value,
            pass: util.id("authPassword").value
        }, data => {
            //console.log(data,"qqq")
            if (data["token"]) {
                //console.log(data,"qqqw")
                localStorage.setItem("token", data["token"])
                util.setToken(data["token"]);

                if (data["access"]) {
                    this.showInterface();
                    document.getElementById('showPaymentInfoBtn').style.display = 'block';

                } else {
                    this.hideInterface()
                }
            } else {
                alert("Ytf")
            }
        })
    };

    this.init();
}

    document.getElementById('logoutBtn').addEventListener('click', function() {
    fetch('/logout', {
        method: 'POST',
        headers: {
            'Authorization': localStorage.getItem('token')
        }
    })
        .then(response => response.text())
        .then(data => {
            console.log(data);
            localStorage.removeItem('token');
            window.location.reload();
        })
        .catch(error => console.error('Ошибка:', error));
});
document.getElementById('deleteUserBtn').addEventListener('click', function() {
    main.deleteUser();
});
document.getElementById('deleteVdsBtn').addEventListener('click', function() {
    main.deleteVds();
});


const main = new Main();
