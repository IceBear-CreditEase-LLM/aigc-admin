import{d as c,r,k as n,l as i,m,Y as g,L as h,au as V}from"./utils-15090c58.js";import{aD as y}from"./index-24ceeda1.js";import{V as w}from"./VFileInput-01190dfb.js";const C=c({__name:"UploadFile",props:{modelValue:{default:""},infos:{default:()=>({})},purpose:{default:""},loading:{type:Boolean,default:!1},showLoading:{type:Boolean,default:!0}},emits:["update:modelValue","update:infos","upload:success","loading"],setup(u,{emit:d}){const l=u,a=d,p=r(),o=r(!1),f=async e=>{if(e.length===0)return;const s={file:e};l.purpose&&(s.purpose=l.purpose),o.value=!0,a("loading",!0);const[_,t]=await V.upload({url:"/files",data:s});t&&(a("update:modelValue",t.fileId),a("update:infos",t),a("upload:success",t)),o.value=!1,a("loading",!1)};return(e,s)=>(n(),i(w,h({ref_key:"refFileInput",ref:p},e.$attrs,{density:"compact",variant:"outlined","hide-details":"auto","onUpdate:modelValue":f,disabled:e.showLoading&&o.value}),{"append-inner":m(()=>[e.showLoading&&o.value?(n(),i(y,{key:0,indeterminate:"",color:"primary",size:20,width:2})):g("",!0)]),_:1},16,["disabled"]))}});export{C as _};
