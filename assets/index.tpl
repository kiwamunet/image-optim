<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <title>image-optim</title>
    <style type="text/css">
      html, body {
        width: 100%;
        height: 100%;
        margin: 0;
        padding: 0;
      }
      html {
        display: table;
      }
      body {
        display: table-cell;
        text-align: center;
      }
      table {
	      border-collapse: collapse;
      }
      td {
	      border: solid 1px;
	      padding: 0.5em;
      }
}
    </style>
</head>
<body>
    <h1>image-optim</h1>
    <div>
    <strong>URL is</strong><br>
    :---------:---------:---------:---------:---------:---------:---------:---------:<br>
    <A Href="{{.URL}}">{{.URL}}</A><br>
    <br>
    </div>
    <div>
    <strong>Commands is</strong><br>
    :---------:---------:---------:---------:---------:---------:---------:---------:<br>
    <script>
      for (var i = 0; i < {{.Commands}}.length; i++){
        document.write({{.Commands}}[i]);
        document.write("<br>");
      }
    // -->
    </script>
    <br>
    </div>
    <div align="center">
    <table>
      <tr>
        <td><img src="data:{{.OriginalMimeType}};base64,{{.OriginalImage}}"></td>
        <td><img src="data:{{.DstMimeType}};base64,{{.DstImage}}"></td>
      </tr>
      <tr>
        <td>Original     {{.OriginalSize}}</td>
        <td>Output     {{.OutputSize}}{{.Compressibility}}</td>
      </tr>
      <tr>
        <td colspan="2">SSIM:{{.SSIM}} PSNR:{{.PSNR}}</td>
      </tr>
    </table>
    </div>
</body>
</html>