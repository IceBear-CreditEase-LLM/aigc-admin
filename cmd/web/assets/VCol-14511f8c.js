import{aV as a,p as y,m as S,a as C,g as k}from"./index-e2ac1ad9.js";import{A as f,c as N,h}from"./utils-41654a3b.js";const i=(()=>a.reduce((e,r)=>(e[r]={type:[Boolean,String,Number],default:!1},e),{}))(),d=(()=>a.reduce((e,r)=>{const t="offset"+f(r);return e[t]={type:[String,Number],default:null},e},{}))(),m=(()=>a.reduce((e,r)=>{const t="order"+f(r);return e[t]={type:[String,Number],default:null},e},{}))(),u={col:Object.keys(i),offset:Object.keys(d),order:Object.keys(m)};function V(e,r,t){let o=e;if(!(t==null||t===!1)){if(r){const s=r.replace(e,"");o+=`-${s}`}return e==="col"&&(o="v-"+o),e==="col"&&(t===""||t===!0)||(o+=`-${t}`),o.toLowerCase()}}const L=["auto","start","end","center","baseline","stretch"],P=y({cols:{type:[Boolean,String,Number],default:!1},...i,offset:{type:[String,Number],default:null},...d,order:{type:[String,Number],default:null},...m,alignSelf:{type:String,default:null,validator:e=>L.includes(e)},...S(),...C()},"VCol"),v=k()({name:"VCol",props:P(),setup(e,r){let{slots:t}=r;const o=N(()=>{const s=[];let l;for(l in u)u[l].forEach(n=>{const b=e[n],c=V(l,n,b);c&&s.push(c)});const g=s.some(n=>n.startsWith("v-col-"));return s.push({"v-col":!g||!e.cols,[`v-col-${e.cols}`]:e.cols,[`offset-${e.offset}`]:e.offset,[`order-${e.order}`]:e.order,[`align-self-${e.alignSelf}`]:e.alignSelf}),s});return()=>{var s;return h(e.tag,{class:[o.value,e.class],style:e.style},(s=t.default)==null?void 0:s.call(t))}}});export{v as V};
