REPL.itというホスティングサービス上にgolangとstripe(決済サービス)を用いたECサイトを作りました。
ホスティングサービスの制限として、データベースサービスとの連携が出来なく、サーバの連続起動時間も一時間ほどなので、
オブジェクト設計にはファサードパターンを採用し、オブジェクトをJson形式でファイルに保存することにしました。

以下がオブジェクト、ルーティングの構成になります。

WalletFacade
  Account: 
    ユーザ様の情報格納パッケージになります。
  Wallet: 
    残高を格納しているパッケージになります。
  Items:
    購入商品情報の格納パッケージになります。
    List - 購入した商品リスト(Done: 決済済みかどうか=bool), 
    SumValue - ユーザ様の決済履歴の合計金額になります。、
    また、パッケージ内にて掲載商品を.envファイルから読み込み、ItemListという変数に格納しています。
  Item:
    商品情報格納パッケージになります。
  ID:
    ユーザ様のstripe_IDを格納する変数になります。

https://facadestripe.hitabacokyou.repl.co/register/:name/:password ユーザ様の登録
https://facadestripe.hitabacokyou.repl.co/charge  決済
https://facadestripe.hitabacokyou.repl.co/login/:name/:password ログイン
https://facadestripe.hitabacokyou.repl.co/logout ログアウト
https://facadestripe.hitabacokyou.repl.co/deposit/:name/:password/:amount 入金
