import{C as _}from"./Confirm-6ae1b937.js";import{V as C}from"./index-e2ac1ad9.js";import{d as p,a9 as u,r as d,U as k,k as y,l as B,m as o,X as i,j as r,_ as a}from"./utils-41654a3b.js";const V=p({__name:"ConfirmByClick",emits:["close","submit"],setup(b,{expose:l,emit:m}){u();const t=d(),s=m,n=()=>{t.value.hide(),s("close")},c=()=>{s("submit",{showLoading:"btn#btnConfirmByClick"})};return l({show({width:e}={}){t.value.show({width:e})},hide(){n()}}),(e,h)=>{const f=k("AiBtn");return y(),B(_,{ref_key:"refConfirm",ref:t,class:"compo-ConfirmByInput"},{title:o(()=>[i(e.$slots,"title")]),text:o(()=>[i(e.$slots,"text")]),buttons:o(()=>[r(C,{size:"small",color:"secondary",variant:"outlined",onClick:n},{default:o(()=>[a("取消")]),_:1}),r(f,{id:"btnConfirmByClick",size:"small",color:"primary",variant:"flat",onClick:c},{default:o(()=>[a("确定")]),_:1})]),_:3},512)}}});export{V as _};
