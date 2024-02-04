import{a3 as B}from"./index-24ceeda1.js";import{d as $,r as z,x as I,z as H,y as V,S as j,k as R,Z as L,Q as F,n as N,W as q,a0 as G,Y as X}from"./utils-15090c58.js";/*! *****************************************************************************
Copyright (c) Microsoft Corporation.

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
PERFORMANCE OF THIS SOFTWARE.
***************************************************************************** */function v(d,t,e,i){return new(e||(e=Promise))(function(s,n){function r(h){try{l(i.next(h))}catch(o){n(o)}}function a(h){try{l(i.throw(h))}catch(o){n(o)}}function l(h){var o;h.done?s(h.value):(o=h.value,o instanceof e?o:new e(function(p){p(o)})).then(r,a)}l((i=i.apply(d,t||[])).next())})}class k{constructor(){this.listeners={}}on(t,e,i){if(this.listeners[t]||(this.listeners[t]=new Set),this.listeners[t].add(e),i==null?void 0:i.once){const s=()=>{this.un(t,s),this.un(t,e)};return this.on(t,s),s}return()=>this.un(t,e)}un(t,e){var i;(i=this.listeners[t])===null||i===void 0||i.delete(e)}once(t,e){return this.on(t,e,{once:!0})}unAll(){this.listeners={}}emit(t,...e){this.listeners[t]&&this.listeners[t].forEach(i=>i(...e))}}const A={decode:function(d,t){return v(this,void 0,void 0,function*(){const e=new AudioContext({sampleRate:t});return e.decodeAudioData(d).finally(()=>e.close())})},createBuffer:function(d,t){return typeof d[0]=="number"&&(d=[d]),function(e){const i=e[0];if(i.some(s=>s>1||s<-1)){const s=i.length;let n=0;for(let r=0;r<s;r++){const a=Math.abs(i[r]);a>n&&(n=a)}for(const r of e)for(let a=0;a<s;a++)r[a]/=n}}(d),{duration:t,length:d[0].length,sampleRate:d[0].length/t,numberOfChannels:d.length,getChannelData:e=>d==null?void 0:d[e],copyFromChannel:AudioBuffer.prototype.copyFromChannel,copyToChannel:AudioBuffer.prototype.copyToChannel}}};function O(d,t){const e=t.xmlns?document.createElementNS(t.xmlns,d):document.createElement(d);for(const[i,s]of Object.entries(t))if(i==="children")for(const[n,r]of Object.entries(t))typeof r=="string"?e.appendChild(document.createTextNode(r)):e.appendChild(O(n,r));else i==="style"?Object.assign(e.style,s):i==="textContent"?e.textContent=s:e.setAttribute(i,s.toString());return e}function _(d,t,e){const i=O(d,t||{});return e==null||e.appendChild(i),i}var U=Object.freeze({__proto__:null,createElement:_,default:_});const Y={fetchBlob:function(d,t,e){return v(this,void 0,void 0,function*(){const i=yield fetch(d,e);if(i.status>=400)throw new Error(`Failed to fetch ${d}: ${i.status} (${i.statusText})`);return function(s,n){v(this,void 0,void 0,function*(){if(!s.body||!s.headers)return;const r=s.body.getReader(),a=Number(s.headers.get("Content-Length"))||0;let l=0;const h=p=>v(this,void 0,void 0,function*(){l+=(p==null?void 0:p.length)||0;const u=Math.round(l/a*100);n(u)}),o=()=>v(this,void 0,void 0,function*(){let p;try{p=yield r.read()}catch{return}p.done||(h(p.value),yield o())});o()})}(i.clone(),t),i.blob()})}};class Q extends k{constructor(t){super(),this.isExternalMedia=!1,t.media?(this.media=t.media,this.isExternalMedia=!0):this.media=document.createElement("audio"),t.mediaControls&&(this.media.controls=!0),t.autoplay&&(this.media.autoplay=!0),t.playbackRate!=null&&this.onceMediaEvent("canplay",()=>{t.playbackRate!=null&&(this.media.playbackRate=t.playbackRate)})}onMediaEvent(t,e,i){return this.media.addEventListener(t,e,i),()=>this.media.removeEventListener(t,e)}onceMediaEvent(t,e){return this.onMediaEvent(t,e,{once:!0})}getSrc(){return this.media.currentSrc||this.media.src||""}revokeSrc(){const t=this.getSrc();t.startsWith("blob:")&&URL.revokeObjectURL(t)}canPlayType(t){return this.media.canPlayType(t)!==""}setSrc(t,e){if(this.getSrc()===t)return;this.revokeSrc();const i=e instanceof Blob&&this.canPlayType(e.type)?URL.createObjectURL(e):t;this.media.src=i}destroy(){this.media.pause(),this.isExternalMedia||(this.media.remove(),this.revokeSrc(),this.media.src="",this.media.load())}setMediaElement(t){this.media=t}play(){return v(this,void 0,void 0,function*(){if(this.media.src)return this.media.play()})}pause(){this.media.pause()}isPlaying(){return!this.media.paused&&!this.media.ended}setTime(t){this.media.currentTime=t}getDuration(){return this.media.duration}getCurrentTime(){return this.media.currentTime}getVolume(){return this.media.volume}setVolume(t){this.media.volume=t}getMuted(){return this.media.muted}setMuted(t){this.media.muted=t}getPlaybackRate(){return this.media.playbackRate}isSeeking(){return this.media.seeking}setPlaybackRate(t,e){e!=null&&(this.media.preservesPitch=e),this.media.playbackRate=t}getMediaElement(){return this.media}setSinkId(t){return this.media.setSinkId(t)}}class T extends k{constructor(t,e){super(),this.timeouts=[],this.isScrollable=!1,this.audioData=null,this.resizeObserver=null,this.lastContainerWidth=0,this.isDragging=!1,this.options=t;const i=this.parentFromOptionsContainer(t.container);this.parent=i;const[s,n]=this.initHtml();i.appendChild(s),this.container=s,this.scrollContainer=n.querySelector(".scroll"),this.wrapper=n.querySelector(".wrapper"),this.canvasWrapper=n.querySelector(".canvases"),this.progressWrapper=n.querySelector(".progress"),this.cursor=n.querySelector(".cursor"),e&&n.appendChild(e),this.initEvents()}parentFromOptionsContainer(t){let e;if(typeof t=="string"?e=document.querySelector(t):t instanceof HTMLElement&&(e=t),!e)throw new Error("Container not found");return e}initEvents(){const t=i=>{const s=this.wrapper.getBoundingClientRect(),n=i.clientX-s.left,r=i.clientX-s.left;return[n/s.width,r/s.height]};this.wrapper.addEventListener("click",i=>{const[s,n]=t(i);this.emit("click",s,n)}),this.wrapper.addEventListener("dblclick",i=>{const[s,n]=t(i);this.emit("dblclick",s,n)}),this.options.dragToSeek&&this.initDrag(),this.scrollContainer.addEventListener("scroll",()=>{const{scrollLeft:i,scrollWidth:s,clientWidth:n}=this.scrollContainer,r=i/s,a=(i+n)/s;this.emit("scroll",r,a)});const e=this.createDelay(100);this.resizeObserver=new ResizeObserver(()=>{e().then(()=>this.onContainerResize()).catch(()=>{})}),this.resizeObserver.observe(this.scrollContainer)}onContainerResize(){const t=this.parent.clientWidth;t===this.lastContainerWidth&&this.options.height!=="auto"||(this.lastContainerWidth=t,this.reRender())}initDrag(){(function(t,e,i,s,n=3,r=0){if(!t)return()=>{};let a=()=>{};const l=h=>{if(h.button!==r)return;h.preventDefault(),h.stopPropagation();let o=h.clientX,p=h.clientY,u=!1;const f=c=>{c.preventDefault(),c.stopPropagation();const w=c.clientX,C=c.clientY,x=w-o,E=C-p;if(u||Math.abs(x)>n||Math.abs(E)>n){const M=t.getBoundingClientRect(),{left:S,top:D}=M;u||(i==null||i(o-S,p-D),u=!0),e(x,E,w-S,C-D),o=w,p=C}},m=()=>{u&&(s==null||s()),a()},y=c=>{c.relatedTarget&&c.relatedTarget!==document.documentElement||m()},b=c=>{u&&(c.stopPropagation(),c.preventDefault())},g=c=>{u&&c.preventDefault()};document.addEventListener("pointermove",f),document.addEventListener("pointerup",m),document.addEventListener("pointerout",y),document.addEventListener("pointercancel",y),document.addEventListener("touchmove",g,{passive:!1}),document.addEventListener("click",b,{capture:!0}),a=()=>{document.removeEventListener("pointermove",f),document.removeEventListener("pointerup",m),document.removeEventListener("pointerout",y),document.removeEventListener("pointercancel",y),document.removeEventListener("touchmove",g),setTimeout(()=>{document.removeEventListener("click",b,{capture:!0})},10)}};t.addEventListener("pointerdown",l)})(this.wrapper,(t,e,i)=>{this.emit("drag",Math.max(0,Math.min(1,i/this.wrapper.getBoundingClientRect().width)))},()=>this.isDragging=!0,()=>this.isDragging=!1)}getHeight(t){return t==null?128:isNaN(Number(t))?t==="auto"&&this.parent.clientHeight||128:Number(t)}initHtml(){const t=document.createElement("div"),e=t.attachShadow({mode:"open"});return e.innerHTML=`
      <style>
        :host {
          user-select: none;
          min-width: 1px;
        }
        :host audio {
          display: block;
          width: 100%;
        }
        :host .scroll {
          overflow-x: auto;
          overflow-y: hidden;
          width: 100%;
          position: relative;
        }
        :host .noScrollbar {
          scrollbar-color: transparent;
          scrollbar-width: none;
        }
        :host .noScrollbar::-webkit-scrollbar {
          display: none;
          -webkit-appearance: none;
        }
        :host .wrapper {
          position: relative;
          overflow: visible;
          z-index: 2;
        }
        :host .canvases {
          min-height: ${this.getHeight(this.options.height)}px;
        }
        :host .canvases > div {
          position: relative;
        }
        :host canvas {
          display: block;
          position: absolute;
          top: 0;
          image-rendering: pixelated;
        }
        :host .progress {
          pointer-events: none;
          position: absolute;
          z-index: 2;
          top: 0;
          left: 0;
          width: 0;
          height: 100%;
          overflow: hidden;
        }
        :host .progress > div {
          position: relative;
        }
        :host .cursor {
          pointer-events: none;
          position: absolute;
          z-index: 5;
          top: 0;
          left: 0;
          height: 100%;
          border-radius: 2px;
        }
      </style>

      <div class="scroll" part="scroll">
        <div class="wrapper" part="wrapper">
          <div class="canvases"></div>
          <div class="progress" part="progress"></div>
          <div class="cursor" part="cursor"></div>
        </div>
      </div>
    `,[t,e]}setOptions(t){if(this.options.container!==t.container){const e=this.parentFromOptionsContainer(t.container);e.appendChild(this.container),this.parent=e}t.dragToSeek&&!this.options.dragToSeek&&this.initDrag(),this.options=t,this.reRender()}getWrapper(){return this.wrapper}getScroll(){return this.scrollContainer.scrollLeft}destroy(){var t;this.container.remove(),(t=this.resizeObserver)===null||t===void 0||t.disconnect()}createDelay(t=10){let e,i;const s=()=>{e&&clearTimeout(e),i&&i()};return this.timeouts.push(s),()=>new Promise((n,r)=>{s(),i=r,e=setTimeout(()=>{e=void 0,i=void 0,n()},t)})}convertColorValues(t){if(!Array.isArray(t))return t||"";if(t.length<2)return t[0]||"";const e=document.createElement("canvas"),i=e.getContext("2d"),s=e.height*(window.devicePixelRatio||1),n=i.createLinearGradient(0,0,0,s),r=1/(t.length-1);return t.forEach((a,l)=>{const h=l*r;n.addColorStop(h,a)}),n}renderBarWaveform(t,e,i,s){const n=t[0],r=t[1]||t[0],a=n.length,{width:l,height:h}=i.canvas,o=h/2,p=window.devicePixelRatio||1,u=e.barWidth?e.barWidth*p:1,f=e.barGap?e.barGap*p:e.barWidth?u/2:0,m=e.barRadius||0,y=l/(u+f)/a,b=m&&"roundRect"in i?"roundRect":"rect";i.beginPath();let g=0,c=0,w=0;for(let C=0;C<=a;C++){const x=Math.round(C*y);if(x>g){const S=Math.round(c*o*s),D=S+Math.round(w*o*s)||1;let W=o-S;e.barAlign==="top"?W=0:e.barAlign==="bottom"&&(W=h-D),i[b](g*(u+f),W,u,D,m),g=x,c=0,w=0}const E=Math.abs(n[C]||0),M=Math.abs(r[C]||0);E>c&&(c=E),M>w&&(w=M)}i.fill(),i.closePath()}renderLineWaveform(t,e,i,s){const n=r=>{const a=t[r]||t[0],l=a.length,{height:h}=i.canvas,o=h/2,p=i.canvas.width/l;i.moveTo(0,o);let u=0,f=0;for(let m=0;m<=l;m++){const y=Math.round(m*p);if(y>u){const g=o+(Math.round(f*o*s)||1)*(r===0?-1:1);i.lineTo(u,g),u=y,f=0}const b=Math.abs(a[m]||0);b>f&&(f=b)}i.lineTo(u,o)};i.beginPath(),n(0),n(1),i.fill(),i.closePath()}renderWaveform(t,e,i){if(i.fillStyle=this.convertColorValues(e.waveColor),e.renderFunction)return void e.renderFunction(t,i);let s=e.barHeight||1;if(e.normalize){const n=Array.from(t[0]).reduce((r,a)=>Math.max(r,Math.abs(a)),0);s=n?1/n:1}e.barWidth||e.barGap||e.barAlign?this.renderBarWaveform(t,e,i,s):this.renderLineWaveform(t,e,i,s)}renderSingleCanvas(t,e,i,s,n,r,a,l){const h=window.devicePixelRatio||1,o=document.createElement("canvas"),p=t[0].length;o.width=Math.round(i*(r-n)/p),o.height=s*h,o.style.width=`${Math.floor(o.width/h)}px`,o.style.height=`${s}px`,o.style.left=`${Math.floor(n*i/h/p)}px`,a.appendChild(o);const u=o.getContext("2d");if(this.renderWaveform(t.map(f=>f.slice(n,r)),e,u),o.width>0&&o.height>0){const f=o.cloneNode(),m=f.getContext("2d");m.drawImage(o,0,0),m.globalCompositeOperation="source-in",m.fillStyle=this.convertColorValues(e.progressColor),m.fillRect(0,0,o.width,o.height),l.appendChild(f)}}renderChannel(t,e,i){return v(this,void 0,void 0,function*(){const s=document.createElement("div"),n=this.getHeight(e.height);s.style.height=`${n}px`,this.canvasWrapper.style.minHeight=`${n}px`,this.canvasWrapper.appendChild(s);const r=s.cloneNode();this.progressWrapper.appendChild(r);const a=t[0].length,l=(g,c)=>{this.renderSingleCanvas(t,e,i,n,Math.max(0,g),Math.min(c,a),s,r)};if(!this.isScrollable)return void l(0,a);const{scrollLeft:h,scrollWidth:o,clientWidth:p}=this.scrollContainer,u=a/o;let f=Math.min(T.MAX_CANVAS_WIDTH,p);if(e.barWidth||e.barGap){const g=e.barWidth||.5,c=g+(e.barGap||g/2);f%c!=0&&(f=Math.floor(f/c)*c)}const m=Math.floor(Math.abs(h)*u),y=Math.floor(m+f*u),b=y-m;l(m,y),yield Promise.all([(()=>v(this,void 0,void 0,function*(){if(m===0)return;const g=this.createDelay();for(let c=m;c>=0;c-=b)yield g(),l(Math.max(0,c-b),c)}))(),(()=>v(this,void 0,void 0,function*(){if(y===a)return;const g=this.createDelay();for(let c=y;c<a;c+=b)yield g(),l(c,Math.min(a,c+b))}))()])})}render(t){return v(this,void 0,void 0,function*(){this.timeouts.forEach(a=>a()),this.timeouts=[],this.canvasWrapper.innerHTML="",this.progressWrapper.innerHTML="",this.options.width!=null&&(this.scrollContainer.style.width=typeof this.options.width=="number"?`${this.options.width}px`:this.options.width);const e=window.devicePixelRatio||1,i=this.scrollContainer.clientWidth,s=Math.ceil(t.duration*(this.options.minPxPerSec||0));this.isScrollable=s>i;const n=this.options.fillParent&&!this.isScrollable,r=(n?i:s)*e;this.wrapper.style.width=n?"100%":`${s}px`,this.scrollContainer.style.overflowX=this.isScrollable?"auto":"hidden",this.scrollContainer.classList.toggle("noScrollbar",!!this.options.hideScrollbar),this.cursor.style.backgroundColor=`${this.options.cursorColor||this.options.progressColor}`,this.cursor.style.width=`${this.options.cursorWidth}px`,this.audioData=t,this.emit("render");try{if(this.options.splitChannels)yield Promise.all(Array.from({length:t.numberOfChannels}).map((a,l)=>{var h;const o=Object.assign(Object.assign({},this.options),(h=this.options.splitChannels)===null||h===void 0?void 0:h[l]);return this.renderChannel([t.getChannelData(l)],o,r)}));else{const a=[t.getChannelData(0)];t.numberOfChannels>1&&a.push(t.getChannelData(1)),yield this.renderChannel(a,this.options,r)}}catch{return}this.emit("rendered")})}reRender(){if(!this.audioData)return;const{scrollWidth:t}=this.scrollContainer,e=this.progressWrapper.clientWidth;if(this.render(this.audioData),this.isScrollable&&t!==this.scrollContainer.scrollWidth){const i=this.progressWrapper.clientWidth;this.scrollContainer.scrollLeft+=i-e}}zoom(t){this.options.minPxPerSec=t,this.reRender()}scrollIntoView(t,e=!1){const{scrollLeft:i,scrollWidth:s,clientWidth:n}=this.scrollContainer,r=t*s,a=i,l=i+n,h=n/2;if(this.isDragging)r+30>l?this.scrollContainer.scrollLeft+=30:r-30<a&&(this.scrollContainer.scrollLeft-=30);else{(r<a||r>l)&&(this.scrollContainer.scrollLeft=r-(this.options.autoCenter?h:0));const o=r-i-h;e&&this.options.autoCenter&&o>0&&(this.scrollContainer.scrollLeft+=Math.min(o,10))}{const o=this.scrollContainer.scrollLeft,p=o/s,u=(o+n)/s;this.emit("scroll",p,u)}}renderProgress(t,e){if(isNaN(t))return;const i=100*t;this.canvasWrapper.style.clipPath=`polygon(${i}% 0, 100% 0, 100% 100%, ${i}% 100%)`,this.progressWrapper.style.width=`${i}%`,this.cursor.style.left=`${i}%`,this.cursor.style.transform=`translateX(-${Math.round(i)===100?this.options.cursorWidth:0}px)`,this.isScrollable&&this.options.autoScroll&&this.scrollIntoView(t,e)}exportImage(t,e,i){return v(this,void 0,void 0,function*(){const s=this.canvasWrapper.querySelectorAll("canvas");if(!s.length)throw new Error("No waveform data");if(i==="dataURL"){const n=Array.from(s).map(r=>r.toDataURL(t,e));return Promise.resolve(n)}return Promise.all(Array.from(s).map(n=>new Promise((r,a)=>{n.toBlob(l=>{l?r(l):a(new Error("Could not export image"))},t,e)})))})}}T.MAX_CANVAS_WIDTH=4e3;class Z extends k{constructor(){super(...arguments),this.unsubscribe=()=>{}}start(){this.unsubscribe=this.on("tick",()=>{requestAnimationFrame(()=>{this.emit("tick")})}),this.emit("tick")}stop(){this.unsubscribe()}destroy(){this.unsubscribe()}}class J extends k{constructor(t=new AudioContext){super(),this.bufferNode=null,this.autoplay=!1,this.playStartTime=0,this.playedDuration=0,this._muted=!1,this.buffer=null,this.currentSrc="",this.paused=!0,this.crossOrigin=null,this.addEventListener=this.on,this.removeEventListener=this.un,this.audioContext=t,this.gainNode=this.audioContext.createGain(),this.gainNode.connect(this.audioContext.destination)}load(){return v(this,void 0,void 0,function*(){})}get src(){return this.currentSrc}set src(t){if(this.currentSrc=t,!t)return this.buffer=null,void this.emit("emptied");fetch(t).then(e=>e.arrayBuffer()).then(e=>this.currentSrc!==t?null:this.audioContext.decodeAudioData(e)).then(e=>{this.currentSrc===t&&(this.buffer=e,this.emit("loadedmetadata"),this.emit("canplay"),this.autoplay&&this.play())})}_play(){var t;this.paused&&(this.paused=!1,(t=this.bufferNode)===null||t===void 0||t.disconnect(),this.bufferNode=this.audioContext.createBufferSource(),this.bufferNode.buffer=this.buffer,this.bufferNode.connect(this.gainNode),this.playedDuration>=this.duration&&(this.playedDuration=0),this.bufferNode.start(this.audioContext.currentTime,this.playedDuration),this.playStartTime=this.audioContext.currentTime,this.bufferNode.onended=()=>{this.currentTime>=this.duration&&(this.pause(),this.emit("ended"))})}_pause(){var t;this.paused||(this.paused=!0,(t=this.bufferNode)===null||t===void 0||t.stop(),this.playedDuration+=this.audioContext.currentTime-this.playStartTime)}play(){return v(this,void 0,void 0,function*(){this._play(),this.emit("play")})}pause(){this._pause(),this.emit("pause")}stopAt(t){var e,i;const s=t-this.currentTime;(e=this.bufferNode)===null||e===void 0||e.stop(this.audioContext.currentTime+s),(i=this.bufferNode)===null||i===void 0||i.addEventListener("ended",()=>{this.bufferNode=null,this.pause()},{once:!0})}setSinkId(t){return v(this,void 0,void 0,function*(){return this.audioContext.setSinkId(t)})}get playbackRate(){var t,e;return(e=(t=this.bufferNode)===null||t===void 0?void 0:t.playbackRate.value)!==null&&e!==void 0?e:1}set playbackRate(t){this.bufferNode&&(this.bufferNode.playbackRate.value=t)}get currentTime(){return this.paused?this.playedDuration:this.playedDuration+this.audioContext.currentTime-this.playStartTime}set currentTime(t){this.emit("seeking"),this.paused?this.playedDuration=t:(this._pause(),this.playedDuration=t,this._play()),this.emit("timeupdate")}get duration(){var t;return((t=this.buffer)===null||t===void 0?void 0:t.duration)||0}get volume(){return this.gainNode.gain.value}set volume(t){this.gainNode.gain.value=t,this.emit("volumechange")}get muted(){return this._muted}set muted(t){this._muted!==t&&(this._muted=t,this._muted?this.gainNode.disconnect():this.gainNode.connect(this.audioContext.destination))}canPlayType(t){return/^(audio|video)\//.test(t)}getGainNode(){return this.gainNode}getChannelData(){const t=[];if(!this.buffer)return t;const e=this.buffer.numberOfChannels;for(let i=0;i<e;i++)t.push(this.buffer.getChannelData(i));return t}}const K={waveColor:"#999",progressColor:"#555",cursorWidth:1,minPxPerSec:0,fillParent:!0,interact:!0,dragToSeek:!1,autoScroll:!0,autoCenter:!0,sampleRate:8e3};class P extends Q{static create(t){return new P(t)}constructor(t){const e=t.media||(t.backend==="WebAudio"?new J:void 0);super({media:e,mediaControls:t.mediaControls,autoplay:t.autoplay,playbackRate:t.audioRate}),this.plugins=[],this.decodedData=null,this.subscriptions=[],this.mediaSubscriptions=[],this.options=Object.assign({},K,t),this.timer=new Z;const i=e?void 0:this.getMediaElement();this.renderer=new T(this.options,i),this.initPlayerEvents(),this.initRendererEvents(),this.initTimerEvents(),this.initPlugins(),Promise.resolve().then(()=>{this.emit("init");const s=this.options.url||this.getSrc()||"";(s||this.options.peaks&&this.options.duration)&&this.load(s,this.options.peaks,this.options.duration)})}updateProgress(t=this.getCurrentTime()){return this.renderer.renderProgress(t/this.getDuration(),this.isPlaying()),t}initTimerEvents(){this.subscriptions.push(this.timer.on("tick",()=>{if(!this.isSeeking()){const t=this.updateProgress();this.emit("timeupdate",t),this.emit("audioprocess",t)}}))}initPlayerEvents(){this.isPlaying()&&(this.emit("play"),this.timer.start()),this.mediaSubscriptions.push(this.onMediaEvent("timeupdate",()=>{const t=this.updateProgress();this.emit("timeupdate",t)}),this.onMediaEvent("play",()=>{this.emit("play"),this.timer.start()}),this.onMediaEvent("pause",()=>{this.emit("pause"),this.timer.stop()}),this.onMediaEvent("emptied",()=>{this.timer.stop()}),this.onMediaEvent("ended",()=>{this.emit("finish")}),this.onMediaEvent("seeking",()=>{this.emit("seeking",this.getCurrentTime())}))}initRendererEvents(){this.subscriptions.push(this.renderer.on("click",(t,e)=>{this.options.interact&&(this.seekTo(t),this.emit("interaction",t*this.getDuration()),this.emit("click",t,e))}),this.renderer.on("dblclick",(t,e)=>{this.emit("dblclick",t,e)}),this.renderer.on("scroll",(t,e)=>{const i=this.getDuration();this.emit("scroll",t*i,e*i)}),this.renderer.on("render",()=>{this.emit("redraw")}),this.renderer.on("rendered",()=>{this.emit("redrawcomplete")}));{let t;this.subscriptions.push(this.renderer.on("drag",e=>{this.options.interact&&(this.renderer.renderProgress(e),clearTimeout(t),t=setTimeout(()=>{this.seekTo(e)},this.isPlaying()?0:200),this.emit("interaction",e*this.getDuration()),this.emit("drag",e))}))}}initPlugins(){var t;!((t=this.options.plugins)===null||t===void 0)&&t.length&&this.options.plugins.forEach(e=>{this.registerPlugin(e)})}unsubscribePlayerEvents(){this.mediaSubscriptions.forEach(t=>t()),this.mediaSubscriptions=[]}setOptions(t){this.options=Object.assign({},this.options,t),this.renderer.setOptions(this.options),t.audioRate&&this.setPlaybackRate(t.audioRate),t.mediaControls!=null&&(this.getMediaElement().controls=t.mediaControls)}registerPlugin(t){return t._init(this),this.plugins.push(t),this.subscriptions.push(t.once("destroy",()=>{this.plugins=this.plugins.filter(e=>e!==t)})),t}getWrapper(){return this.renderer.getWrapper()}getScroll(){return this.renderer.getScroll()}getActivePlugins(){return this.plugins}loadAudio(t,e,i,s){return v(this,void 0,void 0,function*(){if(this.emit("load",t),!this.options.media&&this.isPlaying()&&this.pause(),this.decodedData=null,!e&&!i){const r=a=>this.emit("loading",a);e=yield Y.fetchBlob(t,r,this.options.fetchParams)}this.setSrc(t,e);const n=s||this.getDuration()||(yield new Promise(r=>{this.onceMediaEvent("loadedmetadata",()=>r(this.getDuration()))}));if(i)this.decodedData=A.createBuffer(i,n||0);else if(e){const r=yield e.arrayBuffer();this.decodedData=yield A.decode(r,this.options.sampleRate)}this.decodedData&&(this.emit("decode",this.getDuration()),this.renderer.render(this.decodedData)),this.emit("ready",this.getDuration())})}load(t,e,i){return v(this,void 0,void 0,function*(){yield this.loadAudio(t,void 0,e,i)})}loadBlob(t,e,i){return v(this,void 0,void 0,function*(){yield this.loadAudio("blob",t,e,i)})}zoom(t){if(!this.decodedData)throw new Error("No audio loaded");this.renderer.zoom(t),this.emit("zoom",t)}getDecodedData(){return this.decodedData}exportPeaks({channels:t=2,maxLength:e=8e3,precision:i=1e4}={}){if(!this.decodedData)throw new Error("The audio has not been decoded yet");const s=Math.min(t,this.decodedData.numberOfChannels),n=[];for(let r=0;r<s;r++){const a=this.decodedData.getChannelData(r),l=[],h=Math.round(a.length/e);for(let o=0;o<e;o++){const p=a.slice(o*h,(o+1)*h);let u=0;for(let f=0;f<p.length;f++){const m=p[f];Math.abs(m)>Math.abs(u)&&(u=m)}l.push(Math.round(u*i)/i)}n.push(l)}return n}getDuration(){let t=super.getDuration()||0;return t!==0&&t!==1/0||!this.decodedData||(t=this.decodedData.duration),t}toggleInteraction(t){this.options.interact=t}setTime(t){super.setTime(t),this.updateProgress(t)}seekTo(t){const e=this.getDuration()*t;this.setTime(e)}playPause(){return v(this,void 0,void 0,function*(){return this.isPlaying()?this.pause():this.play()})}stop(){this.pause(),this.setTime(0)}skip(t){this.setTime(this.getCurrentTime()+t)}empty(){this.load("",[[0]],.001)}setMediaElement(t){this.unsubscribePlayerEvents(),super.setMediaElement(t),this.initPlayerEvents()}exportImage(t="image/png",e=1,i="dataURL"){return v(this,void 0,void 0,function*(){return this.renderer.exportImage(t,e,i)})}destroy(){this.emit("destroy"),this.plugins.forEach(t=>t.destroy()),this.subscriptions.forEach(t=>t()),this.unsubscribePlayerEvents(),this.timer.destroy(),this.renderer.destroy(),super.destroy()}}P.BasePlugin=class extends k{constructor(d){super(),this.subscriptions=[],this.options=d}onInit(){}_init(d){this.wavesurfer=d,this.onInit()}destroy(){this.emit("destroy"),this.subscriptions.forEach(d=>d())}},P.dom=U;const tt={class:"compo-aiAudio"},et=["src"],it=["element-loading-text"],rt=$({__name:"AiAudio",props:{src:{default:""},gender:{default:3},type:{default:"simple"},time:{}},setup(d){const t=z(),e=I({style:{isReadyComplex:!1,classNameComplex:"opacity-0"},formData:{}}),{style:i,formData:s}=H(e),n=d,r=()=>{let{type:l,gender:h,src:o}=n;if(l=="complex"&&o){B(t.value).html("");let p={1:{waveColor:"#38bdf8",progressColor:"#38bdf850",cursorColor:"#CCC"},2:{waveColor:"#f980e9",progressColor:"#f980e950",cursorColor:"#ccc"},3:{waveColor:"#475569",progressColor:"#47556950",cursorColor:"#ccc"}},{waveColor:u,progressColor:f,cursorColor:m}=p[h];P.create({container:t.value,waveColor:u,progressColor:f,cursorColor:m,cursorWidth:3,mediaControls:!0,url:o,autoplay:!1,interact:!0}).on("ready",()=>{a()})}},a=()=>{e.style.classNameComplex="w-75 opacity-0",setTimeout(()=>{e.style.classNameComplex="w-100 opacity-0"},100),setTimeout(()=>{e.style.isReadyComplex=!0,e.style.classNameComplex=""},200)};return V(()=>{let{type:l,src:h}=n;l=="complex"&&h&&setTimeout(()=>{r()},100)}),(l,h)=>{const o=j("loading");return R(),L("div",tt,[l.type=="simple"?(R(),L("audio",{key:0,style:{"vertical-align":"top"},src:l.src,controls:""}," Your browser does not support the audio element. ",8,et)):l.type=="complex"?F((R(),L("div",{key:1,class:"aiAudio-complex","element-loading-text":l.$t("aiAudio.loadingText")},[q("div",{ref_key:"refBoxComplex",ref:t,class:G(N(i).classNameComplex)},null,2)],8,it)),[[o,!N(i).isReadyComplex]]):X("",!0)])}}});export{rt as _};
