<script>
  import { onMount } from 'svelte';
  import { url, goto } from '@roxi/routify';
  import ws from '../../../utils/socket';

  export let gameid;

  let baseURL = window.location.protocol + '//' + window.location.host;
  let headerdiv = document.getElementsByClassName('header')[0];


  // Hide navbar in this page
  headerdiv.style.display = 'none';

  onMount(() => {
    // init websocket
    const socket = ws.init(gameid, 'Scanpage');

    let qrcodelink = baseURL + $url(`/${gameid}/game`);
    new QRCode(document.getElementById('qrcode'), {
      text: qrcodelink,
      width: 256,
      height: 256,
      correctLevel: QRCode.CorrectLevel.H,
    });
    // Apply style to qrcode here
    let qrcodeimg = document.querySelector('div#qrcode img');
    qrcodeimg.style.borderRadius = '10px';

    socket.addEventListener("redirect", () => {
      $goto($url(`/${gameid}/scoreboard`))
    })
  });
</script>

<div class="flex flex-col space-y-10 p-6">
  <div>
    <h1 class="text-center text-6xl underline font-semibold">
      Start page for game with ID
      <i>{gameid}</i>
    </h1>
  </div>
  <div>
    <h3 class="text-center text-4xl">Setup a new game like so:</h3>
  </div>
  <div class="mx-auto">
    <ul class="list-disc">
      <li>Connect to the local wifi</li>
      <li>Scan the QR code with your smartphone or tablet</li>
      <li>Setup your desired game</li>
      <li>Start the game</li>
    </ul>
  </div>
  <div>
    <p class="text-center">
      If you cannot scan the QRCode navigate your mobile browser to:
      <a
        class="font-bold"
        href={$url(`/${gameid}/game`)}>{baseURL + $url(`/${gameid}/game`)}</a>
    </p>
  </div>
  <div class="mx-auto" id="qrcode" />
</div>
