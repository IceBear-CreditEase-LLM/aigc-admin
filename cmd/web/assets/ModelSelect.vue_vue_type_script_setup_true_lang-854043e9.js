import{az as p,b3 as c}from"./index-24ceeda1.js";import{d as V,r as v,o as b,k as j,l as k,V as y,m as B,X as M,L as O,n as S,f as _,au as g}from"./utils-15090c58.js";const o="modelName",w=V({__name:"ModelSelect",props:{modelValue:{default:null},defaultFirstValue:{type:Boolean,default:!1},returnObject:{type:Boolean,default:!1}},emits:["update:modelValue"],setup(u,{emit:r}){const n=u,m=r,a=v([]),l=p(n,"modelValue",m),d=async()=>{const[e,t]=await g.get({url:"/channels/models"});t&&(a.value=t.list,i())},i=()=>{const{modelValue:e,defaultFirstValue:t,returnObject:s}=n;e?s&&(l.value=a.value.find(f=>f[o]===e.modelName)||e):t&&(s?l.value=a.value[0]:l.value=a.value[0][o])};return b(()=>{d()}),(e,t)=>(j(),k(c,O({modelValue:S(l),"onUpdate:modelValue":t[0]||(t[0]=s=>_(l)?l.value=s:null)},e.$attrs,{placeholder:"请选择模型",items:a.value,"item-title":o,"item-value":o,variant:"outlined","return-object":n.returnObject}),y({_:2},[e.$slots.prepend?{name:"prepend",fn:B(()=>[M(e.$slots,"prepend")]),key:"0"}:void 0]),1040,["modelValue","items","return-object"]))}});export{w as _};