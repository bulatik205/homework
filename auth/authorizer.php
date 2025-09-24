<?php
require_once __DIR__ . '/../API/config.php';

$authCookie = $_COOKIE["cookie_auth"];
$authValid = false;

if (isset($_COOKIE['auth_hash'])) {
    $stmt = $main_pdo->prepare("SELECT * FROM users WHERE auth_hash = ?");
    $stmt->execute([$_COOKIE['auth_hash']]);
    $data = $stmt->fetch(PDO::FETCH_ASSOC);
    if (!empty($data)) {
        $authValid = true;
    }
}
?>