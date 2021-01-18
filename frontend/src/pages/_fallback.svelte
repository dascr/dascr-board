<script>
  let baseURL = `${window.location.protocol}//${window.location.host}`;
</script>

<div class="flex flex-col p-10 space-y-10">
  <div>
    <h1 class="text-center font-bold underline text-5xl">
      Ooops... You missed a shot!
    </h1>
  </div>
  <div>
    <p>The page you are looking for was not found.</p>
    <p>
      How about you start at:
      <a class="font-bold" href={baseURL}>{baseURL}</a>
    </p>
    <p>Or you use the menu (top left) to navigate.</p>
  </div>
</div>
