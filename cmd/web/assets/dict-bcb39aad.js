import{d as M,r as d,x as P,c as z,U as w,k as L,l as A,m as l,j as e,Y as Z,al as G,am as J,W as u,P as X,av as R,a1 as ee,$ as S,ak as re,Z as Q,_ as g,ax as ne,ac as ie,o as de,S as ce,aA as E,Q as ue,n as O,aB as pe,F as me}from"./utils-766f017d.js";import{_ as fe}from"./BaseBreadcrumb.vue_vue_type_style_index_0_lang-43b6f3c0.js";import{_ as H}from"./UiParentCard.vue_vue_type_script_setup_true_lang-3544285d.js";import{A as _e}from"./AlertBlock-6eaee3a4.js";import{a5 as B,_ as te,V as ye}from"./index-88a5d611.js";import{V as ve}from"./VTextarea-25f7dce9.js";import{V as be}from"./VForm-3cc00baa.js";import{_ as he}from"./ConfirmByInput.vue_vue_type_style_index_0_lang-ff70adad.js";import{_ as ge}from"./ConfirmByClick.vue_vue_type_style_index_0_lang-65078dc2.js";import{V as T}from"./VCol-07cdd7be.js";import{V as le}from"./VRow-5c8abc48.js";import"./IconInfoCircle-fd84073d.js";import"./VAlert-8363e7da.js";import"./Confirm-583e63df.js";const U=v=>(G("data-v-7128444a"),v=v(),J(),v),xe=U(()=>u("label",{class:"required"},"字典编号",-1)),ke=U(()=>u("label",{class:"required"},"字典名称",-1)),De=U(()=>u("label",{class:"required"},"字典键值",-1)),Ve=U(()=>u("label",{class:"required"},"字典排序",-1)),we=U(()=>u("label",{class:"required"},"值类型",-1)),Ce=U(()=>u("label",null,"字典备注",-1)),$e=M({__name:"DictForm",props:{type:{default:"add"},parentId:{default:0}},setup(v,{expose:F}){const C={parentId:0,code:"",dictLabel:"",dictValue:"",dictType:null,sort:0,remark:""},k=v,o=d({...C}),f=d(),n=P({code:[t=>!!t||"请输入字典编号"],dictLabel:[t=>!!t||"请输入字典名称"],dictValue:[t=>!!t||"请输入字典键值"],sort:[t=>t!==""||"请输入字典排序"],dictType:[t=>!!t||"请选择值类型"]}),x=z(()=>k.type==="edit"),m=z(()=>x.value||k.type==="add"&&k.parentId!==0);return F({reset:(t={})=>{o.value={...C,...t}},setFormData:t=>{o.value={...t}},getFormData(){return{...o.value,parentId:k.parentId}},getRef(){return f.value}}),(t,i)=>{const b=w("Select");return L(),A(be,{ref_key:"refForm",ref:f,class:"my-form"},{default:l(()=>[e(B,{type:"text",placeholder:"请输入字典编号","hide-details":"auto",clearable:"",rules:n.code,modelValue:o.value.code,"onUpdate:modelValue":i[0]||(i[0]=p=>o.value.code=p),disabled:m.value},{prepend:l(()=>[xe]),_:1},8,["rules","modelValue","disabled"]),e(B,{type:"text",placeholder:"请输入字典名称","hide-details":"auto",clearable:"",rules:n.dictLabel,modelValue:o.value.dictLabel,"onUpdate:modelValue":i[1]||(i[1]=p=>o.value.dictLabel=p)},{prepend:l(()=>[ke]),_:1},8,["rules","modelValue"]),t.parentId!==0?(L(),A(B,{key:0,type:"text",placeholder:"请输入字典键值","hide-details":"auto",clearable:"",rules:n.dictValue,modelValue:o.value.dictValue,"onUpdate:modelValue":i[2]||(i[2]=p=>o.value.dictValue=p)},{prepend:l(()=>[De]),_:1},8,["rules","modelValue"])):Z("",!0),e(B,{type:"number",placeholder:"请输入字典排序","hide-details":"auto",rules:n.sort,modelValue:o.value.sort,"onUpdate:modelValue":i[3]||(i[3]=p=>o.value.sort=p),modelModifiers:{number:!0}},{prepend:l(()=>[Ve]),_:1},8,["rules","modelValue"]),e(b,{placeholder:"请选择值类型",rules:n.dictType,mapDictionary:{code:"sys_dict_type"},modelValue:o.value.dictType,"onUpdate:modelValue":i[4]||(i[4]=p=>o.value.dictType=p)},{prepend:l(()=>[we]),_:1},8,["rules","modelValue"]),e(ve,{modelValue:o.value.remark,"onUpdate:modelValue":i[5]||(i[5]=p=>o.value.remark=p),modelModifiers:{trim:!0},placeholder:"请输入字典备注",clearable:""},{prepend:l(()=>[Ce]),_:1},8,["modelValue"])]),_:1},512)}}});const ae=te($e,[["__scopeId","data-v-7128444a"]]),Ie={class:"mx-auto mt-3",style:{width:"540px"}},Fe=M({__name:"CreateDictPane",emits:["submit"],setup(v,{expose:F,emit:C}){const k=C,o=P({operateType:"add"}),f=d(),n=d(),x=async({showLoading:m})=>{const{valid:r}=await n.value.getRef().validate();if(r){const y=n.value.getFormData(),t={url:"",method:""};o.operateType=="add"?(t.url="/sys/dict",t.method="post"):(t.url=`/sys/dict/${y.id}`,t.method="put");const[i,b]=await R[t.method]({showLoading:m,showSuccess:!0,url:t.url,data:y});b&&(f.value.hide(),k("submit"))}else ee.warning("请处理页面标错的地方后，再尝试提交")};return F({show({title:m,operateType:r,infos:y}){f.value.show({title:m}),o.operateType=r,X(()=>{o.operateType==="add"?n.value.reset():n.value.setFormData(y)})}}),(m,r)=>{const y=w("Pane");return L(),A(y,{ref_key:"refPane",ref:f,onSubmit:x},{default:l(()=>[u("div",Ie,[e(ae,{ref_key:"refDictForm",ref:n},null,512)])]),_:1},512)}}}),Te=v=>(G("data-v-7706e3d6"),v=v(),J(),v),Pe={class:"d-flex justify-space-between align-center w-100 pr-2 overflow-hidden"},Le={class:"text-truncate"},Se=["onClick"],Be=["onClick"],Re=Te(()=>u("br",null,null,-1)),Ue={class:"text-primary font-weight-black"},qe=M({__name:"ConfigDictPane",emits:["refresh"],setup(v,{expose:F,emit:C}){const k=C,o=d([]),f=d(),n=d(),x=d(),m=d(),r=P({type:"edit",parentId:0,title:"编辑"}),y=d(0),t=d(null),i=d(),b=P({id:null,currentLabel:""}),p=d(!1),$=async()=>{const[a,c]=await R.get({url:"/sys/dict",data:{parentId:y.value}});c&&(o.value=c.list||[],t.value||q(c.list[0]),X(()=>{I()}))},I=()=>{n.value.setCurrentKey(t.value)},q=a=>{r.title=`编辑（${a.dictLabel}）`,r.type="edit",t.value=a.id;const c={id:a.id,parentId:a.parentId,code:a.code,dictValue:a.dictValue,dictLabel:a.dictLabel,dictType:a.dictType,sort:a.sort,remark:a.remark};r.parentId=a.parentId,m.value.setFormData(c)},K=async()=>{const{valid:a}=await m.value.getRef().validate();if(a){const c=m.value.getFormData(),s={url:"",method:""};r.type=="add"?(s.url="/sys/dict",s.method="post"):(s.url=`/sys/dict/${c.id}`,s.method="put");const[j,W]=await R[s.method]({showLoading:x.value.$el,showSuccess:!0,url:s.url,data:c});W&&(await $(),r.type==="add"&&m.value.reset({code:c.code}),c.parentId===0&&r.type==="edit"&&(p.value=!0))}else ee.warning("请处理页面标错的地方后，再尝试提交")},_=a=>{r.title=`新增子项(${a.dictLabel})`,r.type="add",r.parentId=a.id,t.value=a.id,I(),m.value.reset({code:a.code}),V()},h=a=>{b.currentLabel=a.dictLabel,b.id=a.id,i.value.show({width:"450px",confirmText:b.currentLabel})},D=async(a={})=>{const[c,s]=await R.delete({...a,showSuccess:!0,url:`/sys/dict/${b.id}`});s&&(i.value.hide(),b.id===t.value&&(t.value=null),$())},V=()=>{x.value&&ne(x.value.$el,[{transformOrigin:"center"},{transform:"scale(1.03)"},{transformOrigin:"center"}],{duration:150,easing:"cubic-bezier(0.4, 0, 0.2, 1)"})};function N(){t.value=null,p.value&&k("refresh")}return F({show({title:a,id:c}){f.value.show({width:900,showActions:!1,title:a}),y.value=c,p.value=!1,$()}}),(a,c)=>{const s=w("el-tree"),j=w("AiBtn"),W=w("Pane");return L(),A(W,{ref_key:"refPane",ref:f,onClose:N},{default:l(()=>[e(le,null,{default:l(()=>[e(T,{cols:"5"},{default:l(()=>[e(H,null,{default:l(()=>[e(s,{ref_key:"refTree",ref:n,data:o.value,"node-key":"id","default-expand-all":"","highlight-current":!0,"expand-on-click-node":!1,props:{label:"dictLabel"},onNodeClick:q},{default:l(({node:oe,data:Y})=>[u("span",Pe,[u("span",Le,S(oe.label),1),u("span",{class:"ml-2",onClick:c[0]||(c[0]=re(()=>{},["stop"]))},[u("a",{class:"link text-info",onClick:se=>_(Y)},"添加",8,Se),Y.parentId!==0?(L(),Q("a",{key:0,class:"link text-error ml-2",onClick:se=>h(Y)},"删除",8,Be)):Z("",!0)])])]),_:1},8,["data"])]),_:1})]),_:1}),e(T,{cols:"7"},{default:l(()=>[e(H,{class:"dict-card",ref_key:"refDictFormCard",ref:x,title:r.title},{action:l(()=>[e(j,{class:"ml-2",size:"small",color:"primary",variant:"flat",onClick:K},{default:l(()=>[g("提交")]),_:1})]),default:l(()=>[e(ae,{ref_key:"refDictForm",ref:m,type:r.type,parentId:r.parentId},null,8,["type","parentId"])]),_:1},8,["title"])]),_:1})]),_:1}),e(ge,{ref_key:"refConfirmDelete",ref:i,onSubmit:D},{text:l(()=>[g(" 这是进行一项操作时必须了解的重要信息"),Re,g(" 您将要删除 "),u("span",Ue,S(b.currentLabel),1),g(" ，确定要继续吗？ ")]),_:1},512)]),_:1},512)}}});const Ae=te(qe,[["__scopeId","data-v-7706e3d6"]]),Me=u("span",{class:"text-primary"},"删除",-1),Ke=u("br",null,null,-1),Ne={class:"text-primary font-weight-black"},We=u("br",null,null,-1),at=M({__name:"dict",setup(v){const{getLabels:F,loadDictTree:C}=ie();C(["sys_dict_type"]);const k=d({title:"系统字典"}),o=d([{text:"系统管理",disabled:!1,href:"#"},{text:"系统字典",disabled:!0,href:"#"}]),f=P({code:"",label:""}),n=P({list:[],total:0}),x=d(),m=d(),r=d(),y=d(),t=P({id:null,currentCode:""}),i=_=>{let h=[];return h.push({text:"字典配置",color:"info",click(){K(_)}}),h.push({text:"删除",color:"error",click(){b(_)}}),h},b=_=>{t.currentCode=_.code,t.id=_.id,y.value.show({width:"450px",confirmText:t.currentCode})},p=async(_={})=>{const[h,D]=await R.delete({..._,showSuccess:!0,url:`/sys/dict/${t.id}`});D&&(y.value.hide(),$())},$=async(_={})=>{const[h,D]=await R.get({url:"/sys/dict",showLoading:r.value.el,data:{...f,..._}});D?(n.list=D.list||[],n.total=D.total):(n.list=[],n.total=0)},I=()=>{r.value.query({page:1})},q=()=>{x.value.show({title:"添加字典",operateType:"add"})},K=_=>{m.value.show({title:`字典配置（${_.dictLabel}）`,id:_.id})};return de(()=>{$()}),(_,h)=>{const D=w("ButtonsInForm"),V=w("el-table-column"),N=w("ButtonsInTable"),a=w("TableWithPager"),c=ce("copy");return L(),Q(me,null,[e(fe,{title:k.value.title,breadcrumbs:o.value},null,8,["title","breadcrumbs"]),e(H,null,{default:l(()=>[e(le,null,{default:l(()=>[e(T,{cols:"12",lg:"3",md:"4",sm:"6"},{default:l(()=>[e(B,{modelValue:f.code,"onUpdate:modelValue":h[0]||(h[0]=s=>f.code=s),label:"请输入字典编号","hide-details":"",clearable:"",onKeyup:E(I,["enter"]),"onClick:clear":I},null,8,["modelValue","onKeyup"])]),_:1}),e(T,{cols:"12",lg:"3",md:"4",sm:"6"},{default:l(()=>[e(B,{modelValue:f.label,"onUpdate:modelValue":h[1]||(h[1]=s=>f.label=s),label:"请输入字典名称","hide-details":"",clearable:"",onKeyup:E(I,["enter"]),"onClick:clear":I},null,8,["modelValue","onKeyup"])]),_:1}),e(T,{cols:"12",lg:"3",md:"4",sm:"6"},{default:l(()=>[e(D,null,{default:l(()=>[e(ye,{color:"primary",onClick:q},{default:l(()=>[g("添加字典")]),_:1})]),_:1})]),_:1}),e(T,{cols:"12"},{default:l(()=>[e(_e,null,{default:l(()=>[g("修改之后将实时生效，请谨慎操作！")]),_:1})]),_:1}),e(T,{cols:"12"},{default:l(()=>[e(a,{onQuery:$,ref_key:"tableWithPagerRef",ref:r,infos:n},{default:l(()=>[e(V,{label:"字典编号",width:"200px",prop:"code"},{default:l(({row:s})=>[ue((L(),Q("span",null,[g(S(s.code),1)])),[[c,s.code]])]),_:1}),e(V,{label:"字典名称",prop:"dictLabel",width:"150px"}),e(V,{label:"字典排序",prop:"sort",width:"100px"}),e(V,{label:"字典类型",width:"100px"},{default:l(({row:s})=>[u("span",null,S(O(F)([["sys_dict_type",s.dictType]])),1)]),_:1}),e(V,{label:"备注",prop:"remark","min-width":"200px"}),e(V,{label:"更新时间","min-width":"160px"},{default:l(({row:s})=>[g(S(O(pe).dateFormat(s.updatedAt,"YYYY-MM-DD HH:mm:ss")),1)]),_:1}),e(V,{label:"操作",width:"160px",fixed:"right"},{default:l(({row:s})=>[e(N,{buttons:i(s)},null,8,["buttons"])]),_:1})]),_:1},8,["infos"])]),_:1})]),_:1})]),_:1}),e(he,{ref_key:"refConfirmDelete",ref:y,onSubmit:p},{text:l(()=>[g(" 此操作将会"),Me,g("正在使用的字典"),Ke,g(" 字典编号："),u("span",Ne,S(t.currentCode),1),We,g(" 你还要继续吗？ ")]),_:1},512),e(Fe,{ref_key:"createDictPaneRef",ref:x,onSubmit:I},null,512),e(Ae,{ref_key:"configDictPaneRef",ref:m,onRefresh:$},null,512)],64)}}});export{at as default};
