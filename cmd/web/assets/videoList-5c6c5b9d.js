import{d as ie,at as re,ac as ce,r as i,x as z,o as ue,U as v,S as de,k as f,Z as x,j as e,m as t,n as _,aA as fe,_ as n,a0 as b,F as L,a5 as N,l as k,W as s,$ as c,ak as me,Q as j,aB as _e,O as pe,a3 as he,av as y}from"./utils-41654a3b.js";import{_ as ge}from"./BaseBreadcrumb.vue_vue_type_style_index_0_lang-cda68b72.js";import{_ as A}from"./UiParentCard.vue_vue_type_script_setup_true_lang-59b0c9d2.js";import{A as ye}from"./AlertBlock-abe84057.js";import{_ as F}from"./ConfirmByClick.vue_vue_type_style_index_0_lang-6c985436.js";import{_ as ve}from"./TaskOverview.vue_vue_type_script_setup_true_lang-acbf50ac.js";import{c as xe,a5 as be,V as Q,a0 as ke,M as we,e as Ce,as as Ve,K as De,I as $e,i as Se,f as Be,at as Ie}from"./index-e2ac1ad9.js";import{_ as Le}from"./DialogLog.vue_vue_type_style_index_0_lang-62c65c8f.js";import{t as Fe,d as T}from"./map-680f8f44.js";import{_ as Me}from"./TagCorner.vue_vue_type_style_index_0_lang-bfe5450d.js";import{V as K}from"./VRow-5318fa5b.js";import{V as p}from"./VCol-14511f8c.js";import{I as Pe,a as Re,b as q}from"./IconEye-77d24a2b.js";import"./IconInfoCircle-205d0c55.js";import"./VAlert-315a0f3d.js";import"./Confirm-6ae1b937.js";import"./TextLog-d251b87b.js";import"./IconLoader-8634a3a9.js";import"./IconCircleCheckFilled-8db1f68c.js";var ze=xe("bolt","IconBolt",[["path",{d:"M13 3l0 7l6 0l-8 11l0 -7l-6 0l8 -11",key:"svg-0"}]]);const Ne=["onClick"],je={class:"pa-3"},Ae={class:"text-h5"},Qe={class:"text-subtitle-1 mt-1 text-medium-emphasis text-truncate"},Te={class:"text-subtitle-2 mt-1 text-medium-emphasis text-truncate"},Ke={class:"d-flex align-center justify-space-between mt-2",style:{height:"32px"}},qe={class:"flex-1-1 d-flex justify-space-between text-medium-emphasis"},He=s("span",{class:"font-weight-bold"},"这是进行一项操作时必须了解的重要信息",-1),Oe=s("br",null,null,-1),Ue=s("span",{class:"text-primary"},"删除",-1),Ee=s("br",null,null,-1),We={class:"text-primary"},Ze=s("br",null,null,-1),Ge=s("span",{class:"font-weight-bold"},"这是进行一项操作时必须了解的重要信息",-1),Je=s("br",null,null,-1),Xe=s("span",{class:"text-primary"},"取消",-1),Ye=s("br",null,null,-1),et={class:"text-primary"},tt=s("br",null,null,-1),st=s("span",{class:"font-weight-bold"},"这是进行一项操作时必须了解的重要信息",-1),at=s("br",null,null,-1),lt=s("span",{class:"text-primary"},"优先合成",-1),ot=s("br",null,null,-1),nt={class:"text-primary"},it=s("br",null,null,-1),St=ie({name:"DigitalVideoList",__name:"videoList",setup(rt){const M=re(),{getLabels:H,loadDictTree:O}=ce(),U=i({title:"视频合成列表"}),E=i([]),h=z({title:"",status:null}),m=i([]),w=i(0),C=i(),V=i(),D=i(),$=i(),S=i(),u=z({uuid:"",title:""}),P=i(),R=i(),B=i(!1);O(["speak_gender"]);const W=async(l={})=>{B.value=!0;const[a,r]=await y.get({url:"/digitalhuman/synthesis/list",showLoading:P.value.$el,data:{...h,...l}});r?(m.value=r.list||[],w.value=r.total):(m.value=[],w.value=0),B.value=!1,R.value.start()},I=()=>{C.value.query({page:1})},g=()=>{C.value.query()},Z=l=>{let a=[];return l!=="waiting"&&a.push({text:"日志",color:"info",icon:Re,key:"log"}),(l==="waiting"||l==="running")&&a.push({text:"取消",color:"error",icon:q,key:"cancel"}),(l==="failed"||l==="cancel"||l==="success")&&a.push({text:"删除",color:"error",icon:q,key:"delete"}),l==="waiting"&&a.push({text:"优速通",color:"info",icon:ze,key:"first"}),a},G=({id:l},a)=>{u.uuid=a.uuid,u.title=a.title,l==="log"?J():l==="cancel"?D.value.show({width:"400px"}):l==="delete"?$.value.show({width:"400px"}):l==="first"&&S.value.show({width:"400px"})},J=async()=>{V.value.show();let[l,a]=await y.get({url:`/api/digitalhuman/synthesis/${u.uuid}/view`});a&&V.value.setContent(a.synthesisLog)},X=async(l={})=>{const[a,r]=await y.put({...l,showSuccess:!0,url:`/api/digitalhuman/synthesis/${u.uuid}/cancel`});r&&(D.value.hide(),g())},Y=async(l={})=>{const[a,r]=await y.delete({...l,showSuccess:!0,url:`/api/digitalhuman/synthesis/${u.uuid}/delete`});r&&($.value.hide(),g())},ee=async(l={})=>{const[a,r]=await y.put({...l,showSuccess:!0,url:`/digitalhuman/synthesis/${u.uuid}/first`});r&&(S.value.hide(),g())},te=()=>{M.push("/digital-human/video-list/edit")},se=({status:l,uuid:a})=>{M.push(`/digital-human/video-list/detail?uuid=${a}`)};return ue(()=>{g()}),(l,a)=>{const r=v("Select"),ae=v("refresh-button"),le=v("ButtonsInForm"),oe=v("NoData"),ne=de("copy");return f(),x(L,null,[e(ge,{title:U.value.title,breadcrumbs:E.value},null,8,["title","breadcrumbs"]),e(A,{class:"mb-3"},{default:t(()=>[e(ve,{config:_(Fe),"request-url":"/digitalhuman/synthesis/count"},null,8,["config"])]),_:1}),e(A,null,{default:t(()=>[e(K,null,{default:t(()=>[e(p,{cols:"12",lg:"3",md:"4",sm:"6"},{default:t(()=>[e(be,{density:"compact",modelValue:h.title,"onUpdate:modelValue":a[0]||(a[0]=o=>h.title=o),label:"请输入标题","hide-details":"",clearable:"",variant:"outlined",onKeyup:fe(I,["enter"]),"onClick:clear":I},null,8,["modelValue","onKeyup"])]),_:1}),e(p,{cols:"12",lg:"3",md:"4",sm:"6"},{default:t(()=>[e(r,{onChange:I,label:"请选择状态",mapDictionary:{code:"digitalhuman_synthesis_status"},modelValue:h.status,"onUpdate:modelValue":a[1]||(a[1]=o=>h.status=o)},null,8,["modelValue"])]),_:1}),e(p,{cols:"12",lg:"3",md:"4",sm:"6"},{default:t(()=>[e(le,null,{default:t(()=>[e(Q,{color:"primary",onClick:te},{default:t(()=>[n("创建视频")]),_:1}),e(ae,{ref_key:"refreshButtonRef",ref:R,onRefresh:g,disabled:B.value},null,8,["disabled"])]),_:1})]),_:1}),e(p,{cols:"12"},{default:t(()=>[e(ye,null,{default:t(()=>[n(" 生成视频的时间可能会比较长，请耐心等待！ ")]),_:1})]),_:1}),e(p,{cols:"12"},{default:t(()=>[e(K,{ref_key:"listContentRef",ref:P,class:b({"justify-center":m.value.length===0})},{default:t(()=>[m.value.length>0?(f(!0),x(L,{key:0},N(m.value,o=>(f(),k(p,{cols:"12",lg:"3",md:"4",sm:"6"},{default:t(()=>[e(ke,{elevation:"10",rounded:"md"},{default:t(()=>[s("a",{class:"card-list-item text-textPrimary text-decoration-none",href:"javascript: void(0)",onClick:d=>se(o)},[e(Me,{class:b(`bg-${_(T)[o.status].color}`)},{default:t(()=>[n(c(_(T)[o.status].text),1)]),_:2},1032,["class"]),e(we,{src:o.digitalHumanPerson.cover,height:"180px",cover:"",class:"rounded-t-md align-end text-right"},{default:t(()=>[s("div",je,[e(Ce,{onClick:a[2]||(a[2]=me(()=>{},["stop"])),class:"bg-surface text-body-2 font-weight-medium",variant:"flat",size:"small",text:`${o.videoDuration||"0s"}/${o.videoSize}`},null,8,["text"])])]),_:2},1032,["src"]),e(Ve,{class:"pa-5"},{default:t(()=>[s("h5",Ae,c(o.title),1),s("p",Qe,c(o.ttsText),1),j((f(),x("p",Te,[n(c(o.uuid),1)])),[[ne,o.uuid]]),s("div",Ke,[s("div",qe,[s("span",null,c(_(_e).dateFromNow(o.createdAt)),1),s("span",null,c(o.digitalHumanPerson.cname)+" ("+c(_(H)([["speak_gender",o.digitalHumanPerson.gender]]))+") ",1)]),e(Q,{class:"ml-6",size:"x-small",color:"inherit",icon:"",variant:"text"},{default:t(()=>[e(_(Pe),{width:"14","stroke-width":"1.5"}),e(De,{activator:"parent"},{default:t(()=>[e($e,{density:"compact","onClick:select":d=>G(d,o)},{default:t(()=>[(f(!0),x(L,null,N(Z(o.status),d=>(f(),k(Se,{key:d.key,value:d.key,"hide-details":"","min-height":"38"},{prepend:t(()=>[(f(),k(pe(d.icon),{size:16,class:b(["mr-2",[`text-${d.color}`]])},null,8,["class"]))]),default:t(()=>[e(Be,{class:b([`text-${d.color}`])},{default:t(()=>[n(c(d.text),1)]),_:2},1032,["class"])]),_:2},1032,["value"]))),128))]),_:2},1032,["onClick:select"])]),_:2},1024)]),_:2},1024)])]),_:2},1024)],8,Ne)]),_:2},1024)]),_:2},1024))),256)):(f(),k(oe,{key:1}))]),_:1},8,["class"]),j(e(Ie,{class:"mt-5",ref_key:"refPager",ref:C,total:w.value,"page-sizes":[12,20,40,60,120],onQuery:W},null,8,["total"]),[[he,m.value.length>0]])]),_:1})]),_:1})]),_:1}),e(Le,{ref_key:"refDialogLog",ref:V},null,512),e(F,{ref_key:"refConfirmDelete",ref:$,onSubmit:Y},{text:t(()=>[He,Oe,n(" 此操作将会"),Ue,n("该视频"),Ee,n(" 标题："),s("span",We,c(u.title),1),Ze,n(" 你还要继续吗？ ")]),_:1},512),e(F,{ref_key:"refConfirmCancel",ref:D,onSubmit:X},{text:t(()=>[Ge,Je,n(" 此操作将会"),Xe,n("该视频"),Ye,n(" 标题："),s("span",et,c(u.title),1),tt,n(" 你还要继续吗？ ")]),_:1},512),e(F,{ref_key:"refConfirmFirst",ref:S,onSubmit:ee},{text:t(()=>[st,at,n(" 此操作将会"),lt,n("该视频"),ot,n(" 标题："),s("span",nt,c(u.title),1),it,n(" 你还要继续吗？ ")]),_:1},512)],64)}}});export{St as default};
