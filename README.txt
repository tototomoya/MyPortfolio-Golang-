REPL.ITというホスティングサービス上にgolangとstripe(決済サービス)を用いたECサイトを作りました。
データベースサービスとの連携が出来なく、一時間ほどでサーバが落ちるので、ファイル上にオブジェクトをJson形式で保存しています。
オブジェクト設計にはファサードパターンを採用しました。

以下がオブジェクト、ルーティングの構成になります。

WalletFacade
  Account: ユーザ情報
  Wallet: 残高
  Items: 購入商品(決済済み=True), SumValueは今までの合計金額になります。
  ID: ユーザのstripe_ID
  
https://facadestripe.hitabacokyou.repl.co/register/:name/:password ユーザ登録
https://facadestripe.hitabacokyou.repl.co/charge  決済
https://facadestripe.hitabacokyou.repl.co/login/:name/:password ログイン
https://facadestripe.hitabacokyou.repl.co/logout ログアウト
https://facadestripe.hitabacokyou.repl.co/deposit/:name/:password/:amount 入金
