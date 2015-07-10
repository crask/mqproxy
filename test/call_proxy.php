<?php

    $ch =curl_init("http://127.0.0.1:9090/produce?format=json");

    $headers = array();
    $headers[] = 'TOPIC:test';
    $headers[] = 'PARTITION_KEY:asasasdasd';

    curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_BINARYTRANSFER, true);

    $data = array(
            'uid'     => 123,
            'uname'   => 'crask',
            'content' => 'welcome to crask',
    );

    $json_str = json_encode($data);
    curl_setopt($ch, CURLOPT_POSTFIELDS, $json_str);

    $t1 = round(microtime(true) * 1000);
    $res = curl_exec($ch);
    $t2 = round(microtime(true) * 1000);
   
    $d = $t2-$t1;
    echo $res . "\ncost: " . $d . "ms\n";
?>
