<?php
$host = 'localhost';
$port = 5434;
$db   = 'portaldata';
$user = 'dev';
$pass = 'devpass';

try {
    $pdo = new PDO("pgsql:host=$host;port=$port;dbname=$db", $user, $pass);
    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
    echo "\nТаблицы в базе $db:\n";
    $tables = $pdo->query("SELECT tablename FROM pg_tables WHERE schemaname='public'")->fetchAll(PDO::FETCH_COLUMN);
    foreach ($tables as $table) {
        echo "\n=== $table ===\n";
        $cols = $pdo->query("SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = '$table'")->fetchAll(PDO::FETCH_ASSOC);
        foreach ($cols as $col) {
            echo $col['column_name'] . "\t" . $col['data_type'] . "\tNULLABLE: " . $col['is_nullable'] . "\n";
        }
    }
} catch (Exception $e) {
    echo "Ошибка: " . $e->getMessage() . "\n";
} 