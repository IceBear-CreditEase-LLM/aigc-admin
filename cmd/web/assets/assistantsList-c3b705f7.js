import{d as W,at as N,r as d,x as v,o as S,U as p,k as Y,Z as j,j as t,m as e,aA as K,_ as n,W as r,$ as m,n as M,aB as Q,F as E,av as y}from"./utils-41654a3b.js";import{_ as H}from"./BaseBreadcrumb.vue_vue_type_style_index_0_lang-cda68b72.js";import{_ as L}from"./UiParentCard.vue_vue_type_script_setup_true_lang-59b0c9d2.js";import{A as U}from"./AlertBlock-abe84057.js";import{_ as q}from"./ConfirmByInput.vue_vue_type_style_index_0_lang-ea20bdb0.js";import{C as Z}from"./CreateAssistantPane-70e8bde2.js";import{V as z}from"./VRow-5318fa5b.js";import{V as f}from"./VCol-14511f8c.js";import{a5 as G,V as J}from"./index-e2ac1ad9.js";import"./IconInfoCircle-205d0c55.js";import"./VAlert-315a0f3d.js";import"./Confirm-6ae1b937.js";import"./ModelSelect.vue_vue_type_script_setup_true_lang-0b7aa11f.js";import"./CustomUpload-79dee0be.js";import"./IconPlus-9cdcb170.js";import"./VTextarea-dd2dd03c.js";import"./VForm-3f7129c5.js";const O=["onClick"],X=["onClick"],tt=r("span",{class:"text-primary font-weight-black"},"删除",-1),et=r("br",null,null,-1),st={class:"text-primary font-weight-black"},at=r("br",null,null,-1),gt=W({__name:"assistantsList",setup(ot){const I=N(),C=d({title:"助手列表"}),V=d([{text:"AI助手",disabled:!1,href:"#"},{text:"助手列表",disabled:!0,href:"#"}]),_=v({name:""}),c=v({list:[],total:0}),h=d(),b=d(),w=d(),u=v({assistantId:""}),B=a=>{let o=[];return o.push({text:"编辑",color:"info",click(){P(a)}}),o.push({text:"删除",color:"error",click(){D(a)}}),o},D=a=>{u.assistantId=a.assistantId,w.value.show({width:"550px",confirmText:u.assistantId})},T=async(a={})=>{const[o,i]=await y.delete({...a,showSuccess:!0,url:`/assistants/${u.assistantId}`});i&&(w.value.hide(),x())},x=async(a={})=>{const[o,i]=await y.get({url:"/assistants/list",showLoading:b.value.el,data:{..._,...a}});i?(c.list=i.list||[],c.total=i.total):(c.list=[],c.total=0)},k=()=>{b.value.query({page:1})},A=()=>{h.value.show({title:"创建助手",operateType:"add"})},P=a=>{h.value.show({title:"编辑助手",infos:a,operateType:"edit"})},g=a=>{I.push(`/ai-assistant/assistants/detail?assistantId=${a}`)};return S(()=>{x()}),(a,o)=>{const i=p("ButtonsInForm"),l=p("el-table-column"),$=p("ButtonsInTable"),F=p("TableWithPager");return Y(),j(E,null,[t(H,{title:C.value.title,breadcrumbs:V.value},null,8,["title","breadcrumbs"]),t(L,null,{default:e(()=>[t(z,null,{default:e(()=>[t(f,{cols:"12",lg:"3",md:"4",sm:"6"},{default:e(()=>[t(G,{modelValue:_.name,"onUpdate:modelValue":o[0]||(o[0]=s=>_.name=s),label:"请输入助手名称","hide-details":"",clearable:"",onKeyup:K(k,["enter"]),"onClick:clear":k},null,8,["modelValue","onKeyup"])]),_:1}),t(f,{cols:"12",lg:"3",md:"4",sm:"6"},{default:e(()=>[t(i,null,{default:e(()=>[t(J,{color:"primary",onClick:A},{default:e(()=>[n("创建助手")]),_:1})]),_:1})]),_:1}),t(f,{cols:"12"},{default:e(()=>[t(U,null,{default:e(()=>[n("修改之后将实时生效，请谨慎操作！")]),_:1})]),_:1}),t(f,{cols:"12"},{default:e(()=>[t(F,{onQuery:x,ref_key:"tableWithPagerRef",ref:b,infos:c},{default:e(()=>[t(l,{label:"助手ID","min-width":"240px"},{default:e(({row:s})=>[r("a",{href:"javascript: void(0)",class:"link",onClick:R=>g(s.assistantId)},m(s.assistantId),9,O)]),_:1}),t(l,{label:"助手名称","min-width":"150px","show-overflow-tooltip":"","class-name":"link-ellipsis-color"},{default:e(({row:s})=>[r("a",{href:"javascript: void(0)",class:"link",onClick:R=>g(s.assistantId)},m(s.name),9,X)]),_:1}),t(l,{label:"工具数量","min-width":"100px"},{default:e(({row:s})=>[n(m(s.tools?s.tools.length:0),1)]),_:1}),t(l,{label:"模型",prop:"modelName","min-width":"200px","show-overflow-tooltip":""}),t(l,{label:"备注",prop:"remark","min-width":"200px"}),t(l,{label:"更新时间","min-width":"165px"},{default:e(({row:s})=>[n(m(M(Q).dateFormat(s.updatedAt,"YYYY-MM-DD HH:mm:ss")),1)]),_:1}),t(l,{label:"操作人",prop:"operator","min-width":"150px","show-overflow-tooltip":""}),t(l,{label:"操作",width:"120px",fixed:"right"},{default:e(({row:s})=>[t($,{buttons:B(s)},null,8,["buttons"])]),_:1})]),_:1},8,["infos"])]),_:1})]),_:1})]),_:1}),t(q,{ref_key:"refConfirmDelete",ref:w,onSubmit:T},{text:e(()=>[n(" 此操作将会"),tt,n("该个人助手，删除之后将无法使用"),et,n(" 助手ID："),r("span",st,m(u.assistantId),1),at,n(" 确定要继续吗？ ")]),_:1},512),t(Z,{ref_key:"createAssistantPaneRef",ref:h,onSubmit:k},null,512)],64)}}});export{gt as default};
