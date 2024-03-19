import{d as T,i as w,r as d,at as A,x as b,o as D,U as m,S as L,k as h,Z as g,j as e,m as t,aA as v,_ as p,Q as R,$ as U,F as W,av as $}from"./utils-41654a3b.js";import{_ as I}from"./BaseBreadcrumb.vue_vue_type_style_index_0_lang-cda68b72.js";import{_ as j}from"./UiParentCard.vue_vue_type_script_setup_true_lang-59b0c9d2.js";import{_ as S}from"./AiAudio.vue_vue_type_style_index_0_lang-a12e7c3a.js";import{V as q}from"./VRow-5318fa5b.js";import{V as i}from"./VCol-14511f8c.js";import{a5 as V,V as x}from"./index-e2ac1ad9.js";const X=T({name:"VoicePrintLibraryList",__name:"index",setup(E){const K=w("provideAspectPage"),C=d({title:"声纹库列表"}),P=d([{text:"智能声纹",disabled:!1,href:"#"},{text:"声纹库列表",disabled:!0,href:"#"}]),f=A(),u=d(),s=b({list:[],total:0}),l=b({userName:"",userKey:""}),o=()=>{u.value.query({page:1})},k=()=>{l.userName="",l.userKey=""},_=async(y={})=>{const[r,n]=await $.get({url:"/voice/list",showLoading:u.value.el,data:{...l,...y}});n?(s.list=n.list||[],s.total=n.total):(s.list=[],s.total=0)},N=()=>{f.push("/voice-print/library-list/register")},F=()=>{f.push("/voice-print/library-list/search")};return K.methods.refreshListPage=()=>{k(),o()},D(()=>{_()}),(y,r)=>{const n=m("ButtonsInForm"),c=m("el-table-column"),B=m("TableWithPager"),Q=L("copy");return h(),g(W,null,[e(I,{title:C.value.title,breadcrumbs:P.value},null,8,["title","breadcrumbs"]),e(j,null,{default:t(()=>[e(q,null,{default:t(()=>[e(i,{cols:"12",lg:"3",md:"4",sm:"6"},{default:t(()=>[e(V,{density:"compact",modelValue:l.userName,"onUpdate:modelValue":r[0]||(r[0]=a=>l.userName=a),label:"用户姓名","hide-details":"",variant:"outlined",clearable:!0,onKeyup:v(o,["enter"]),"onClick:clear":o},null,8,["modelValue","onKeyup"])]),_:1}),e(i,{cols:"12",lg:"3",md:"4",sm:"6"},{default:t(()=>[e(V,{density:"compact",modelValue:l.userKey,"onUpdate:modelValue":r[1]||(r[1]=a=>l.userKey=a),label:"用户标识","hide-details":"",variant:"outlined",clearable:!0,onKeyup:v(o,["enter"]),"onClick:clear":o},null,8,["modelValue","onKeyup"])]),_:1}),e(i,{cols:"12",lg:"3",md:"4",sm:"6"},{default:t(()=>[e(n,null,{default:t(()=>[e(x,{color:"primary",onClick:N},{default:t(()=>[p("声纹注册")]),_:1}),e(x,{color:"secondary",onClick:F},{default:t(()=>[p("声纹查询")]),_:1})]),_:1})]),_:1}),e(i,{cols:"12"},{default:t(()=>[e(B,{onQuery:_,ref_key:"tableWithPagerRef",ref:u,infos:s},{default:t(()=>[e(c,{label:"用户姓名",prop:"userName","min-width":"200px"}),e(c,{label:"用户标识","min-width":"200px"},{default:t(({row:a})=>[R((h(),g("span",null,[p(U(a.userKey),1)])),[[Q,a.userKey]])]),_:1}),e(c,{label:"音频文件","min-width":"300px"},{default:t(({row:a})=>[e(S,{src:a==null?void 0:a.s3Url},null,8,["src"])]),_:1})]),_:1},8,["infos"])]),_:1})]),_:1})]),_:1})],64)}}});export{X as default};
