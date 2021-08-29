<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';

  let gameid = $page.params.gameid;

  let baseURL = window.location.protocol + '//' + window.location.host;

  onMount(async () => {
    // Hide navbar in this page
    let headerdiv = document.getElementsByClassName('header')[0];
    headerdiv.style.display = 'none';

    // init websocket
    const ws = await import('$utils/socket');
    const socket = ws.init(gameid, 'Scanpage');

    let qrcodelink = baseURL + `/${gameid}/game`;
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
      goto(`/${gameid}/scoreboard`);
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
        href={`/${gameid}/game`}>{baseURL + `/${gameid}/game`}</a>
    </p>
  </div>
  <div class="mx-auto" id="qrcode" />
</div>
