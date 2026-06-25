package handlers

import (
	"net/http"
)

func MapHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html class="dark" lang="en">
<head>
<meta charset="utf-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<title>Horror Map - Horror Facts</title>
<script src="https://cdn.tailwindcss.com?plugins=forms,container-queries"></script>
<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@700;900&family=Inter:wght@400&family=JetBrains+Mono:wght@500&display=swap" rel="stylesheet"/>
<link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap" rel="stylesheet"/>
<script>
tailwind.config = {
    darkMode: "class",
    theme: {
        extend: {
            colors: {
                "primary": "#c8c6c6", "background": "#121414", "on-surface": "#ffffff",
                "surface-container": "#1e2020", "surface-container-lowest": "#0c0f0f",
                "outline-variant": "#474747", "on-surface-variant": "#e2e2e2",
                "primary-container": "#474747", "on-primary-container": "#ffffff",
                "secondary": "#c8c6c5",
            },
            fontFamily: {
                "display-lg": ["Playfair Display"], "headline-lg": ["Playfair Display"],
                "headline-md": ["Playfair Display"], "body-md": ["Inter"],
                "body-lg": ["Inter"], "label-md": ["JetBrains Mono"], "caption": ["Inter"],
            }
        }
    }
}
</script>
<style>
body { background-color: #0a0a0a; color: #e2e2e2; margin: 0; min-height: 100vh; display: flex; flex-direction: column; position: relative; overflow-x: hidden; }
body::before { content: ""; position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; pointer-events: none; z-index: 50; opacity: 0.03; background-image: url('data:image/svg+xml,%%3Csvg viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg"%%3E%%3Cfilter id="n"%%3E%%3CfeTurbulence type="fractalNoise" baseFrequency="0.65" numOctaves="3" stitchTiles="stitch"/%%3E%%3C/filter%%3E%%3Crect width="100%%25" height="100%%25" filter="url(%%23n)"/%%3E%%3C/svg%%3E'); }
body::after { content: ""; position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; pointer-events: none; z-index: 49; background: radial-gradient(circle at center, transparent 30%%, rgba(10,10,10,0.8) 100%%); }
.glass-card { background-color: rgba(26,26,26,0.6); backdrop-filter: blur(12px); border: 1px solid rgba(240,240,240,0.1); }
.tag-chip { background-color: #1a1a1a; border: 1px solid #2e2e2e; padding: 4px 8px; font-family: 'JetBrains Mono', monospace; font-size: 12px; color: #fff; }
.map-container { position: relative; width: 100%; height: 500px; border-radius: 8px; overflow: hidden; border: 1px solid rgba(240,240,240,0.1); }
.map-container img { width: 100%; height: 100%; object-fit: cover; filter: brightness(0.4) invert(0.05); }
.map-pin { position: absolute; width: 10px; height: 10px; background: #c8c6c6; border-radius: 50%; cursor: pointer; transform: translate(-50%, -50%); animation: pulse 2s infinite; z-index: 2; }
.map-pin:hover { background: #fff; animation: none; }
@keyframes pulse { 0%, 100% { box-shadow: 0 0 0 0 rgba(200,198,198,0.6); } 50% { box-shadow: 0 0 0 8px rgba(200,198,198,0); } }
.sidebar { transition: transform 0.4s ease-in-out; }
.sidebar.open { transform: translateX(0); }
</style>
</head>
<body>

<nav class="w-full z-50 bg-surface-container/60 backdrop-blur-xl border-b border-outline-variant/20">
    <div class="flex justify-between items-center px-8 py-4 max-w-[1200px] mx-auto">
        <a href="/" class="text-secondary text-2xl font-bold no-underline" style="font-family: 'Playfair Display', serif;">Horror Facts</a>
        <div class="flex gap-6">
            <a href="/" class="text-on-surface-variant hover:text-primary text-sm font-bold uppercase tracking-wider no-underline">Home</a>
            <a href="/movies" class="text-on-surface-variant hover:text-primary text-sm font-bold uppercase tracking-wider no-underline">Archive</a>
            <a href="/map" class="text-primary font-bold border-b-2 border-primary pb-1 text-sm uppercase tracking-wider no-underline">Map</a>
        </div>
    </div>
</nav>

<main class="flex-grow pt-24 relative flex flex-col z-10 px-4 md:px-8 max-w-[1200px] mx-auto w-full pb-16">
    <header class="mb-8 text-center">
        <h1 class="text-5xl md:text-6xl font-black text-primary mb-4" style="font-family: 'Playfair Display', serif;">The Geography of Dread</h1>
        <p class="text-lg text-on-surface-variant max-w-2xl mx-auto opacity-80">Every nightmare has a coordinate. Click the markers to uncover the real stories.</p>
    </header>

    <div class="map-container" id="mapContainer">
        <img src="/static/images/world_map.png" alt="World Map"/>
        <div id="pinsContainer"></div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-8" id="locationCards"></div>
</main>

<div class="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 hidden" id="overlay" onclick="closeSidebar()"></div>

<aside class="sidebar fixed top-0 right-0 h-full w-full md:w-[400px] bg-[#1a1a1a]/95 backdrop-blur-xl border-l border-outline-variant/20 z-50 transform translate-x-full flex flex-col" id="sidebar">
    <div class="p-6 border-b border-outline-variant/20 flex justify-between items-center">
        <h2 class="text-xl font-bold text-primary" style="font-family: 'Playfair Display', serif;" id="sidebarTitle">Location</h2>
        <button onclick="closeSidebar()" class="text-on-surface-variant hover:text-primary"><span class="material-symbols-outlined">close</span></button>
    </div>
    <div class="p-6 flex-grow overflow-y-auto">
        <div class="flex flex-wrap gap-2 mb-4" id="sidebarTags"></div>
        <p class="text-on-surface-variant text-lg leading-relaxed mb-6" id="sidebarFact"></p>
        <h3 class="text-sm font-bold uppercase tracking-wider text-outline-variant mb-3">Related Films</h3>
        <ul class="space-y-2 text-on-surface-variant" id="sidebarFilms"></ul>
    </div>
</aside>

<footer class="w-full py-6 bg-surface-container-lowest border-t border-outline-variant/10 z-10 text-center">
    <p class="text-xs text-on-surface-variant uppercase tracking-wider">EVERY SCREAM HAS A SOURCE.</p>
</footer>

<script>
var locations = [
    {x: 620, y: 165, title: "Moscow", film: "House on the Embankment / White Hand", fact: "Kotelnicheskaya building arrests (1937-38). Metro ghost legend (1970s).", tags: ["RUSSIAN"]},
    {x: 595, y: 155, title: "St. Petersburg", film: "Ghost of Leningrad", fact: "Murders styled after the siege years.", tags: ["RUSSIAN"]},
    {x: 645, y: 175, title: "Samara", film: "Call from the Past", fact: "Calls from a non-existent number predicting deaths.", tags: ["RUSSIAN"]},
    {x: 695, y: 180, title: "Novosibirsk", film: "The Raven", fact: "Prophet case: 7 of 12 predictions came true.", tags: ["RUSSIAN"]},
    {x: 785, y: 150, title: "Oymyakon", film: "Frost on the Skin", fact: "Coldest place in Russia. Murders during abnormal frosts.", tags: ["RUSSIAN"]},
    {x: 570, y: 140, title: "Karelia", film: "Master of the Swamps", fact: "30+ disappearances in Mshinskoye Swamp.", tags: ["RUSSIAN"]},
    {x: 625, y: 170, title: "Kazan", film: "Touch of Darkness", fact: "Illegal LSD experiments in psychiatric hospital.", tags: ["RUSSIAN"]},
    {x: 340, y: 185, title: "Chernobyl", film: "Devil of Chernobyl", fact: "Exclusion zone legends.", tags: ["RUSSIAN"]},
    {x: 715, y: 195, title: "Lake Baikal", film: "Black Baikal", fact: "Sunken NKVD ships.", tags: ["RUSSIAN"]},
    {x: 260, y: 195, title: "Amityville, NY", film: "The Amityville Horror", fact: "112 Ocean Avenue. DeFeo murders.", tags: ["FOREIGN"]},
    {x: 210, y: 180, title: "Plainfield, WI", film: "The Texas Chain Saw Massacre", fact: "Ed Gein - human skin furniture.", tags: ["FOREIGN"]},
    {x: 280, y: 185, title: "Harrisville, RI", film: "The Conjuring", fact: "Perron family haunting.", tags: ["FOREIGN"]},
    {x: 225, y: 210, title: "St. Louis, MO", film: "The Exorcist", fact: "Roland Doe exorcism, 1949.", tags: ["FOREIGN"]},
    {x: 440, y: 170, title: "London, UK", film: "Hellraiser", fact: "1980s subculture inspired Cenobites.", tags: ["FOREIGN"]},
    {x: 890, y: 210, title: "Tokyo, Japan", film: "The Ring / Ju-On", fact: "Cursed videotapes and ju-on houses.", tags: ["FOREIGN"]},
    {x: 460, y: 175, title: "Darmstadt, Germany", film: "The Exorcism of Emily Rose", fact: "Anneliese Michel exorcism, 1976.", tags: ["FOREIGN"]}
];

var pinsContainer = document.getElementById('pinsContainer');

locations.forEach(function(loc, i) {
    var pin = document.createElement('div');
    pin.className = 'map-pin';
    pin.style.left = (loc.x / 1000 * 100) + '%';
    pin.style.top = (loc.y / 500 * 100) + '%';
    pin.title = loc.title;
    pin.onclick = function() { openLocation(i); };
    pinsContainer.appendChild(pin);
});

var cardsContainer = document.getElementById('locationCards');
locations.forEach(function(loc, i) {
    var card = document.createElement('div');
    card.className = 'glass-card rounded-lg p-4 cursor-pointer hover:border-primary transition-all';
    card.onclick = function() { openLocation(i); };
    card.innerHTML = '<h3 class="text-lg font-bold text-primary mb-1" style="font-family: Playfair Display, serif;">' + loc.title + '</h3>' +
        '<p class="text-on-surface-variant text-sm mb-2">' + loc.film + '</p>' +
        '<span class="tag-chip">[' + loc.tags[0] + ']</span>';
    cardsContainer.appendChild(card);
});

function openLocation(i) {
    var d = locations[i];
    document.getElementById('sidebarTitle').textContent = d.title;
    document.getElementById('sidebarFact').textContent = d.fact;
    document.getElementById('sidebarTags').innerHTML = d.tags.map(function(t) { return '<span class="tag-chip">[' + t + ']</span>'; }).join('');
    document.getElementById('sidebarFilms').innerHTML = '<li class="text-on-surface-variant">' + d.film + '</li>';
    document.getElementById('sidebar').classList.add('open');
    document.getElementById('overlay').classList.remove('hidden');
}

function closeSidebar() {
    document.getElementById('sidebar').classList.remove('open');
    document.getElementById('overlay').classList.add('hidden');
}
</script>
</body>
</html>`

	w.Write([]byte(html))
}
