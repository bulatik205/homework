<?php
require_once __DIR__ . "/../auth/authorizer.php";
?>

<?php if ($authValid): ?>
    <!DOCTYPE html>
    <html lang="en">

    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Auth</title>
        <link rel="stylesheet" href="../css/auth/index.css">
        <link rel="stylesheet" href="../css/font/index.css">
    </head>

    <body>
        <header>
            <h1>homework : : :</h1>
        </header>

        <div class="nav">
            <button onclick="navClick(`register`)" class="active-type" id="buttonR">Регистрация</button>
            <button onclick="navClick(`login`)" id="buttonL">Вход</button>
        </div>

        <main>
            <div class="main--" id="reg">
                <div class="main__">
                    <h1>Регистрация</h1>

                    <div class="main--__form">
                        <label for="usernameR">Придумайте имя</label>
                        <input type="text" name="usernameR" id="usernameR" minlength="4" maxlength="255">

                        <label for="passwordR">Придумайте пароль</label>
                        <input type="text" name="passwordR" id="passwordR" minlength="4" maxlength="255">

                        <div class="telegram-auth">
                            <label for="telegramR" class="telegram-t">Код подтверждения</label>
                            <label for="telegramR" class="telegram-p"><a href="https://t.me/bulatik_test_bot">Бот</a> отправит его вам</label>
                            <input type="number" name="telegramR" id="telegramR" placeholder="homework-XXXXXX" style="display: none;">
                            <button onclick="">Отправить код</button>
                            <info id="telegram-infoR">
                                <label id="telegram-info-titleR">Новый код через: 60 секунд</label>
                            </info>
                        </div>

                        <div class="rules">
                            <input type="checkbox" name="rules" id="rules">
                            <label for="rules" id="rulesText">Согласен с <a href="">правилами</a></label>
                        </div>

                        <button onclick="reg()">Регистрация</button>
                    </div>
                </div>
            </div>

            <div class="main--" style="display: none;" id="log">
                <div class="main__">
                    <h1>Вход</h1>

                    <div class="main--log__type">
                        <button onclick="logNavClick(`p`)" id="buttonLP" class="active-log-type">С паролем</button>
                        <button onclick="logNavClick(`b`)" id="buttonLB">Через бота</button>
                    </div>

                    <div class="main--__form" id="p">
                        <label for="usernameL">Введите имя</label>
                        <input type="text" name="usernameL" id="usernameL">

                        <label for="passwordL">Введите пароль</label>
                        <input type="text" name="passwordL" id="passwordL">

                        <button onclick="log(`password`)">Вход</button>
                    </div>

                    <div class="main--__form" style="display: none;" id="b">
                        <div class="telegram-auth">
                            <div class="telegram-auth-warring">
                                <label>Зайдите в <a href="https://t.me/bulatik_test_bot">бота</a>, получите ссылку для авторизации и вернитесь, <label style="color: rgb(120, 170, 170);">перейдя по ней</label></label>
                            </div>

                            <label for="telegram" class="telegram-t">Код подтверждения</label>
                            <label for="telegram" class="telegram-p"><a href="https://t.me/bulatik_test_bot">Бот</a> отправит его вам</label>
                            <input type="number" name="telegram" id="telegramL" placeholder="homework-XXXXXX">
                            <button onclick="">Отправить код</button>
                            <info id="telegram-infoL">
                                <label id="telegram-info-titleL">Новый код через: 60 секунд</label>
                            </info>
                        </div>

                        <button onclick="log(`telegram`)">Вход</button>
                    </div>
                </div>
            </div>
        </main>

        <div class="errors" style="display: none;" id="errors">
            <div>
                <h1 id="error-title"></h1>
                <p id="error-p"></p>
            </div>
        </div>

        <script src="../handlers/auth-handler.js"></script>
        <script src="../handlers/auth-ajax.js"></script>
        <script src="../handlers/auth-validator.js"></script>
    </body>

    </html>

<?php else: ?>
    <?php require("../logger/auth.php") ?>
<?php endif ?>