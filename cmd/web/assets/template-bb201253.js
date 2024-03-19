import{d as W,x as P,r as _,c as z,U as f,k as F,l as H,m as l,W as u,j as e,_ as n,av as q,al as L,am as Z,ac as j,o as G,S as J,Z as $,aA as O,Q as R,$ as T,n as N,aB as X,F as ee}from"./utils-766f017d.js";import{_ as le}from"./BaseBreadcrumb.vue_vue_type_style_index_0_lang-43b6f3c0.js";import{_ as te}from"./UiParentCard.vue_vue_type_script_setup_true_lang-3544285d.js";import{A as ae}from"./AlertBlock-6eaee3a4.js";import{_ as oe}from"./ConfirmByInput.vue_vue_type_style_index_0_lang-ff70adad.js";import{_ as w}from"./Explain.vue_vue_type_style_index_0_lang-395bde7c.js";import{a5 as v,ae as se,_ as ne,V as re}from"./index-88a5d611.js";import{V as de}from"./VSwitch-2a955bcd.js";import{V as ue}from"./VTextarea-25f7dce9.js";import{V as me}from"./VForm-3cc00baa.js";import{C as pe}from"./ChipStatus-4b3e688d.js";import{V as ie}from"./VRow-5c8abc48.js";import{V as k}from"./VCol-07cdd7be.js";import"./IconInfoCircle-fd84073d.js";import"./VAlert-8363e7da.js";import"./Confirm-583e63df.js";const D=x=>(L("data-v-d9604e47"),x=x(),Z(),x),ce={class:"mx-auto mt-3",style:{width:"750px"}},_e=D(()=>u("label",{class:"required"},"模版名称",-1)),fe=D(()=>u("label",{class:"required"},"基础模型",-1)),be=D(()=>u("label",{class:"required"},"模版类型",-1)),he={class:"required"},Ve={class:"required"},ye={class:"required"},ve={class:"required"},xe=D(()=>u("label",{class:"required"},"启用状态",-1)),ge=D(()=>u("label",null,"备注",-1)),Te=W({__name:"CreateTemplatePane",emits:["submit"],setup(x,{expose:I,emit:M}){const S={name:"",baseModel:null,templateType:null,scriptFile:"",content:"",trainImage:"",outputDir:"/data/ft-model/",baseModelPath:"",enabled:!1,remark:""},b=M,p=P({operateType:"add"}),a=_({...S}),h=_(),y=_(),m=P({name:[s=>/^[a-zA-Z0-9-.]+$/.test(s)||"只允许字母、数字、“-” 、“.”"],baseModel:[s=>!!s||"请选择基础模型"],templateType:[s=>!!s||"请选择模版类型"],scriptFile:[s=>!!s||"请输入脚本文件"],trainImage:[s=>!!s||"请输入镜像完整地址"],baseModelPath:[s=>!!s||"请输入基础模型路径"],content:[s=>!!s||"请输入脚本模版"]}),U=z(()=>p.operateType==="edit"),B=async({valid:s,showLoading:o})=>{if(s){const d={url:"",method:""};p.operateType=="add"?(d.url="/sys/template",d.method="post"):(d.url=`/sys/template/${a.value.name}`,d.method="put");const[C,g]=await q[d.method]({showLoading:o,showSuccess:!0,url:d.url,data:a.value});g&&(h.value.hide(),b("submit"))}};return I({show({title:s,operateType:o,infos:d}){h.value.show({title:s,refForm:y}),p.operateType=o,p.operateType==="add"?a.value={...S}:a.value={...d}}}),(s,o)=>{const d=f("Select"),C=f("CodeMirror"),g=f("Pane");return F(),H(g,{ref_key:"refPane",ref:h,onSubmit:B},{default:l(()=>[u("div",ce,[e(me,{ref_key:"refForm",ref:y,class:"my-form"},{default:l(()=>[e(v,{type:"text",placeholder:"只允许字母、数字、“-” 、“.”","hide-details":"auto",clearable:"",rules:m.name,modelValue:a.value.name,"onUpdate:modelValue":o[0]||(o[0]=t=>a.value.name=t),disabled:U.value},{prepend:l(()=>[_e]),_:1},8,["rules","modelValue","disabled"]),e(d,{placeholder:"请选择基础模型",rules:m.baseModel,modelValue:a.value.baseModel,"onUpdate:modelValue":o[1]||(o[1]=t=>a.value.baseModel=t),mapAPI:{url:"/models",data:{pageSize:-1,isPrivate:!0,enabled:!0},labelField:"modelName",valueField:"modelName"},"hide-details":"auto"},{prepend:l(()=>[fe]),_:1},8,["rules","modelValue"]),e(d,{placeholder:"请选择模版类型",rules:m.templateType,mapDictionary:{code:"template_type"},modelValue:a.value.templateType,"onUpdate:modelValue":o[2]||(o[2]=t=>a.value.templateType=t)},{prepend:l(()=>[be]),_:1},8,["rules","modelValue"]),e(v,{type:"text",placeholder:"请输入镜像完整地址","hide-details":"auto",clearable:"",rules:m.trainImage,modelValue:a.value.trainImage,"onUpdate:modelValue":o[3]||(o[3]=t=>a.value.trainImage=t)},{prepend:l(()=>[u("label",he,[n("训练镜像 "),e(w,null,{default:l(()=>[n("请提前将Docker镜像制作好并上传到镜像仓库")]),_:1})])]),_:1},8,["rules","modelValue"]),e(v,{type:"text",placeholder:"请输入模型在容器的绝对路径","hide-details":"auto",clearable:"",rules:m.baseModelPath,modelValue:a.value.baseModelPath,"onUpdate:modelValue":o[4]||(o[4]=t=>a.value.baseModelPath=t)},{prepend:l(()=>[u("label",Ve,[n("基础模型路径 "),e(w,null,{default:l(()=>[n("请输入模型所存储的路径")]),_:1})])]),_:1},8,["rules","modelValue"]),e(v,{type:"text",placeholder:"/data/ft-model/","hide-details":"auto",clearable:"",modelValue:a.value.outputDir,"onUpdate:modelValue":o[5]||(o[5]=t=>a.value.outputDir=t)},{prepend:l(()=>[u("label",null,[n("输出目录 "),e(w,null,{default:l(()=>[n("模型训练所保存的目录")]),_:1})])]),_:1},8,["modelValue"]),e(v,{type:"text",placeholder:"/app/train.sh","hide-details":"auto",clearable:"",rules:m.scriptFile,modelValue:a.value.scriptFile,"onUpdate:modelValue":o[6]||(o[6]=t=>a.value.scriptFile=t)},{prepend:l(()=>[u("label",ye,[n("脚本文件 "),e(w,null,{default:l(()=>[n("训练脚本文件")]),_:1})])]),_:1},8,["rules","modelValue"]),e(se,{"hide-details":"auto",rules:m.content,modelValue:a.value.content,"onUpdate:modelValue":o[8]||(o[8]=t=>a.value.content=t),"center-affix":!1},{prepend:l(()=>[u("label",ve,[n("脚本模版内容 "),e(w,null,{default:l(()=>[n("脚本模版，通常为启动训练脚本的Shell")]),_:1})])]),default:l(()=>[e(C,{modelValue:a.value.content,"onUpdate:modelValue":o[7]||(o[7]=t=>a.value.content=t),language:"shell",placeholder:"请输入脚本模版"},null,8,["modelValue"])]),_:1},8,["rules","modelValue"]),e(de,{modelValue:a.value.enabled,"onUpdate:modelValue":o[9]||(o[9]=t=>a.value.enabled=t),color:"primary","hide-details":"auto"},{prepend:l(()=>[xe]),_:1},8,["modelValue"]),e(ue,{modelValue:a.value.remark,"onUpdate:modelValue":o[10]||(o[10]=t=>a.value.remark=t),modelModifiers:{trim:!0},placeholder:"请输入备注"},{prepend:l(()=>[ge]),_:1},8,["modelValue"])]),_:1},512)])]),_:1},512)}}});const we=ne(Te,[["__scopeId","data-v-d9604e47"]]),ke={class:"text-primary font-weight-black"},Pe=u("br",null,null,-1),Ee=W({__name:"template",setup(x){const{getLabels:I}=j(),M=_({title:"模版管理"}),S=_([{text:"系统管理",disabled:!1,href:"#"},{text:"模版管理",disabled:!0,href:"#"}]),b=P({name:"",templateType:null}),p=P({list:[],total:0}),a=_(),h=_(),y=_(),m=P({name:""}),U=t=>{let i=[];return i.push({text:"删除",color:"error",click(){B(t)}}),i.push({text:"编辑",color:"info",click(){g(t)}}),i},B=t=>{m.name=t.name,y.value.show({width:"450px",confirmText:m.name})},s=async(t={})=>{const[i,V]=await q.delete({...t,showSuccess:!0,url:`/sys/template/${m.name}`});V&&(y.value.hide(),o())},o=async(t={})=>{const[i,V]=await q.get({url:"/sys/template",showLoading:h.value.el,data:{...b,...t}});V?(p.list=V.list||[],p.total=V.total):(p.list=[],p.total=0)},d=()=>{h.value.query({page:1})},C=()=>{a.value.show({title:"添加模版",operateType:"add"})},g=t=>{a.value.show({title:"编辑模版",infos:t,operateType:"edit"})};return G(()=>{o()}),(t,i)=>{const V=f("Select"),Q=f("ButtonsInForm"),c=f("el-table-column"),Y=f("ButtonsInTable"),E=f("TableWithPager"),A=J("copy");return F(),$(ee,null,[e(le,{title:M.value.title,breadcrumbs:S.value},null,8,["title","breadcrumbs"]),e(te,null,{default:l(()=>[e(ie,null,{default:l(()=>[e(k,{cols:"12",lg:"3",md:"4",sm:"6"},{default:l(()=>[e(v,{modelValue:b.name,"onUpdate:modelValue":i[0]||(i[0]=r=>b.name=r),label:"请输入模型","hide-details":"",clearable:"",onKeyup:O(d,["enter"]),"onClick:clear":d},null,8,["modelValue","onKeyup"])]),_:1}),e(k,{cols:"12",lg:"3",md:"4",sm:"6"},{default:l(()=>[e(V,{modelValue:b.templateType,"onUpdate:modelValue":i[1]||(i[1]=r=>b.templateType=r),mapDictionary:{code:"template_type"},label:"请选择模版类型","hide-details":"",onChange:d},null,8,["modelValue"])]),_:1}),e(k,{cols:"12",lg:"3",md:"4",sm:"6"},{default:l(()=>[e(Q,null,{default:l(()=>[e(re,{color:"primary",onClick:C},{default:l(()=>[n("添加模版")]),_:1})]),_:1})]),_:1}),e(k,{cols:"12"},{default:l(()=>[e(ae,null,{default:l(()=>[n("修改之后将实时生效，请谨慎操作！")]),_:1})]),_:1}),e(k,{cols:"12"},{default:l(()=>[e(E,{onQuery:o,ref_key:"tableWithPagerRef",ref:h,infos:p},{default:l(()=>[e(c,{label:"名称",width:"200px","show-overflow-tooltip":""},{default:l(({row:r})=>[R((F(),$("span",null,[n(T(r.name),1)])),[[A,r.name]])]),_:1}),e(c,{label:"基础模型",prop:"baseModel",width:"180px"}),e(c,{label:"最长上下文",prop:"maxTokens",width:"110px"}),e(c,{label:"模版类型",width:"100px"},{default:l(({row:r})=>[n(T(N(I)([["template_type",r.templateType]])),1)]),_:1}),e(c,{label:"镜像",width:"200px","show-overflow-tooltip":""},{default:l(({row:r})=>[R((F(),$("span",null,[n(T(r.trainImage),1)])),[[A,r.trainImage]])]),_:1}),e(c,{label:"状态",width:"100px"},{default:l(({row:r})=>[e(pe,{modelValue:r.enabled,"onUpdate:modelValue":K=>r.enabled=K},null,8,["modelValue","onUpdate:modelValue"])]),_:1}),e(c,{label:"备注",prop:"remark","min-width":"200px"}),e(c,{label:"更新时间","min-width":"165px"},{default:l(({row:r})=>[n(T(N(X).dateFormat(r.updatedAt,"YYYY-MM-DD HH:mm:ss")),1)]),_:1}),e(c,{label:"操作",width:"120px",fixed:"right"},{default:l(({row:r})=>[e(Y,{buttons:U(r)},null,8,["buttons"])]),_:1})]),_:1},8,["infos"])]),_:1})]),_:1})]),_:1}),e(oe,{ref_key:"refConfirmDelete",ref:y,onSubmit:s},{text:l(()=>[n(" 您将要删除"),u("span",ke,T(m.name),1),n("模版，删除之后将无法使用该模版创建微调任务"),Pe,n(" 确定要继续吗？ ")]),_:1},512),e(we,{ref_key:"createTemplatePaneRef",ref:a,onSubmit:d},null,512)],64)}}});export{Ee as default};
