import{ax as i,x as d,D as _,ay as f,U as v,k as a,l as s,m as C,az as g,n as T,O as h}from"./utils-15090c58.js";import{a3 as c}from"./index-24ceeda1.js";const k={__name:"AspectPage",setup(x){const l=i(),t=d({scrollTop:{}});_("provideAspectPage",t);const r=e=>e.matched[e.matched.length-1].components.default.__name,p=e=>(l.meta.aspectPageInclude||[]).includes(e);return f((e,m,n)=>{let o=r(m),u=r(e);p(o)&&(t.scrollTop[o]=c(document).scrollTop()),p(u)?c(document).scrollTop(t.scrollTop[u]||"0"):c(document).scrollTop(0),n()}),(e,m)=>{const n=v("router-view");return a(),s(n,null,{default:C(({Component:o})=>[(a(),s(g,{include:T(l).meta.aspectPageInclude},[(a(),s(h(o)))],1032,["include"]))]),_:1})}}};export{k as default};