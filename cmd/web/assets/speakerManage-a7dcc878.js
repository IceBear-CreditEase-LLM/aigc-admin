import{j as e,L as G,c as Q,F as L,d as Y,ac as j,x as F,r as x,U as S,k as P,l as N,m as a,W as u,_ as b,Z as q,a5 as le,n as $,au as A,al as te,am as oe,o as se,S as re,av as ne,Q as de,$ as w,aw as ue}from"./utils-15090c58.js";import{_ as ie}from"./BaseBreadcrumb.vue_vue_type_style_index_0_lang-f6213a89.js";import{_ as pe}from"./UiParentCard.vue_vue_type_script_setup_true_lang-0e150107.js";import{A as me}from"./AlertBlock-c06c49ed.js";import{_ as ce}from"./ConfirmByInput.vue_vue_type_style_index_0_lang-77722b42.js";import{_ as fe}from"./AiAudio.vue_vue_type_style_index_0_lang-ab9d5699.js";import{_ as _e}from"./Explain.vue_vue_type_style_index_0_lang-9177768e.js";import{_ as ve}from"./UploadFile.vue_vue_type_script_setup_true_lang-4cfa6d0e.js";import{p as E,ac as be,g as K,b as H,ad as z,ae as Ve,af as ge,ag as ke,a6 as W,ab as ye,A as he,ah as xe,ai as B,a4 as Ie,aj as Se,a5 as M,ak as Ce,J as Pe,M as we,_ as De,V as Ue}from"./index-24ceeda1.js";import{V as Re}from"./VAlert-30650f94.js";import{V as Te}from"./VSwitch-83b6bbc0.js";import{V as $e}from"./VTextarea-b1274484.js";import{V as Fe}from"./VForm-123b6ecf.js";import{V as Ne}from"./VRow-ab8e6ed5.js";import{V as D}from"./VCol-e2fcfa86.js";import"./IconInfoCircle-03c84ef6.js";import"./Confirm-aedf99ab.js";import"./VFileInput-01190dfb.js";const Ge=E({...be({falseIcon:"$radioOff",trueIcon:"$radioOn"})},"VRadio"),qe=K()({name:"VRadio",props:Ge(),setup(n,C){let{slots:y}=C;return H(()=>e(z,G(n,{class:["v-radio",n.class],style:n.style,type:"radio"}),y)),{}}});const Ae=E({height:{type:[Number,String],default:"auto"},...Ve(),...ge(ke(),["multiple"]),trueIcon:{type:W,default:"$radioOn"},falseIcon:{type:W,default:"$radioOff"},type:{type:String,default:"radio"}},"VRadioGroup"),Be=K()({name:"VRadioGroup",inheritAttrs:!1,props:Ae(),emits:{"update:modelValue":n=>!0},setup(n,C){let{attrs:y,slots:h}=C;const U=ye(),_=Q(()=>n.id||`radio-group-${U}`),f=he(n,"modelValue");return H(()=>{const[V,t]=xe(y),m=B.filterProps(n),v=z.filterProps(n),I=h.label?h.label({label:n.label,props:{for:_.value}}):n.label;return e(B,G({class:["v-radio-group",n.class],style:n.style},V,m,{modelValue:f.value,"onUpdate:modelValue":i=>f.value=i,id:_.value}),{...h,default:i=>{let{id:c,messagesId:R,isDisabled:T,isReadonly:d}=i;return e(L,null,[I&&e(Ie,{id:c.value},{default:()=>[I]}),e(Se,G(v,{id:c.value,"aria-describedby":R.value,defaultsTarget:"VRadio",trueIcon:n.trueIcon,falseIcon:n.falseIcon,type:n.type,disabled:T.value,readonly:d.value,"aria-labelledby":I?c.value:void 0,multiple:!1},t,{modelValue:f.value,"onUpdate:modelValue":o=>f.value=o}),h)])}})}),{}}}),k=n=>(te("data-v-184fa553"),n=n(),oe(),n),Me={class:"mx-auto mt-3",style:{width:"540px"}},Le={class:"required"},Oe=k(()=>u("label",{class:"required"},"标识",-1)),We=k(()=>u("label",{class:"required"},"姓名",-1)),Qe=k(()=>u("label",{class:"required"},"性别",-1)),Ye=k(()=>u("label",{class:"required"},"年龄段",-1)),je=k(()=>u("label",{class:"required"},"语言",-1)),Ee=k(()=>u("label",{class:"required"},"风格",-1)),Ke=k(()=>u("label",{class:"required"},"适应范围",-1)),He=k(()=>u("label",null,"头像",-1)),ze=k(()=>u("label",{class:"required"},"启用",-1)),Je=k(()=>u("label",null,"备注",-1)),Ze=Y({__name:"CreateSpeakerPane",emits:["submit"],setup(n,{expose:C,emit:y}){const h={provider:null,speakName:"",speakCname:"",gender:1,ageGroup:null,lang:null,speakStyle:null,area:null,headImgFileId:"",enabled:!1,remark:""},U=y,f=j().options.speak_gender,V=F({operateType:"add"}),t=x({...h}),m=x(null),v=x(),I=x(),i=F({provider:[d=>!!d||"请选择供应商"],speakName:[d=>!!d||"请输入标识"],speakCname:[d=>!!d||"请输入姓名"],ageGroup:[d=>!!d||"请选择年龄段"],lang:[d=>!!d||"请选择语言"],speakStyle:[d=>!!d||"请选择风格"],area:[d=>!!d||"请选择适应范围"]}),c=Q(()=>V.operateType==="edit"),R=()=>{t.value.headImgFileId="",m.value=null},T=async({valid:d,showLoading:o})=>{if(d){const s={url:"",method:""};V.operateType=="add"?(s.url="/voice/speak",s.method="post"):(s.url=`/voice/speak/${t.value.id}`,s.method="put");const[p,l]=await A[s.method]({showLoading:o,showSuccess:!0,url:s.url,data:t.value});l&&(v.value.hide(),U("submit"))}};return C({show({title:d,operateType:o,infos:s}){v.value.show({title:d,refForm:I}),V.operateType=o,V.operateType==="add"?(t.value={...h},m.value=null):(t.value={...s},m.value={s3Url:s.headImg})}}),(d,o)=>{const s=S("Select"),p=S("Pane");return P(),N(p,{ref_key:"refPane",ref:v,onSubmit:T},{default:a(()=>[u("div",Me,[e(Fe,{ref_key:"refForm",ref:I,class:"my-form"},{default:a(()=>[e(s,{placeholder:"请选择供应商",rules:i.provider,mapDictionary:{code:"speak_provider"},modelValue:t.value.provider,"onUpdate:modelValue":o[0]||(o[0]=l=>t.value.provider=l),disabled:c.value},{prepend:a(()=>[u("label",Le,[b("供应 "),e(_e,null,{default:a(()=>[b("供应商指的是外部服务提供，自己有服务请选择Local")]),_:1})])]),_:1},8,["rules","modelValue","disabled"]),e(M,{type:"text",placeholder:"请输入标识","hide-details":"auto",clearable:"",rules:i.speakName,modelValue:t.value.speakName,"onUpdate:modelValue":o[1]||(o[1]=l=>t.value.speakName=l),disabled:c.value},{prepend:a(()=>[Oe]),_:1},8,["rules","modelValue","disabled"]),e(M,{type:"text",placeholder:"请输入姓名","hide-details":"auto",clearable:"",rules:i.speakCname,modelValue:t.value.speakCname,"onUpdate:modelValue":o[2]||(o[2]=l=>t.value.speakCname=l)},{prepend:a(()=>[We]),_:1},8,["rules","modelValue"]),e(Be,{"hide-details":"auto",modelValue:t.value.gender,"onUpdate:modelValue":o[3]||(o[3]=l=>t.value.gender=l),inline:"",disabled:c.value},{prepend:a(()=>[Qe]),default:a(()=>[(P(!0),q(L,null,le($(f),l=>(P(),N(qe,{label:l.label,color:"primary",value:l.value},null,8,["label","value"]))),256))]),_:1},8,["modelValue","disabled"]),e(s,{placeholder:"请选择年龄段",rules:i.ageGroup,mapDictionary:{code:"speak_age_group"},modelValue:t.value.ageGroup,"onUpdate:modelValue":o[4]||(o[4]=l=>t.value.ageGroup=l)},{prepend:a(()=>[Ye]),_:1},8,["rules","modelValue"]),e(s,{placeholder:"请选择语言",rules:i.lang,mapDictionary:{code:"speak_lang"},modelValue:t.value.lang,"onUpdate:modelValue":o[5]||(o[5]=l=>t.value.lang=l),disabled:c.value},{prepend:a(()=>[je]),_:1},8,["rules","modelValue","disabled"]),e(s,{placeholder:"请选择风格",rules:i.speakStyle,mapDictionary:{code:"speak_style"},modelValue:t.value.speakStyle,"onUpdate:modelValue":o[6]||(o[6]=l=>t.value.speakStyle=l)},{prepend:a(()=>[Ee]),_:1},8,["rules","modelValue"]),e(s,{placeholder:"请选择适应范围",rules:i.area,mapDictionary:{code:"speak_area"},modelValue:t.value.area,"onUpdate:modelValue":o[7]||(o[7]=l=>t.value.area=l)},{prepend:a(()=>[Ke]),_:1},8,["rules","modelValue"]),e(B,{"hide-details":"auto"},{prepend:a(()=>[He]),default:a(()=>[m.value&&m.value.s3Url?(P(),N(Re,{key:0,color:"borderColor",variant:"outlined",density:"compact"},{close:a(()=>[e(Ce,{class:"text-24 opacity-50 cursor-pointer",color:"textPrimary",onClick:R},{default:a(()=>[b("mdi-close-circle")]),_:1})]),default:a(()=>[e(Pe,{size:"60"},{default:a(()=>[e(we,{transition:!1,src:m.value.s3Url,alt:"上传成功后的头像",cover:""},null,8,["src"])]),_:1})]),_:1})):(P(),N(ve,{key:1,accept:"image/*",modelValue:t.value.headImgFileId,"onUpdate:modelValue":o[8]||(o[8]=l=>t.value.headImgFileId=l),infos:m.value,"onUpdate:infos":o[9]||(o[9]=l=>m.value=l),"prepend-icon":null,"prepend-inner-icon":"mdi-camera"},null,8,["modelValue","infos"]))]),_:1}),e(Te,{modelValue:t.value.enabled,"onUpdate:modelValue":o[10]||(o[10]=l=>t.value.enabled=l),color:"primary","hide-details":"auto"},{prepend:a(()=>[ze]),_:1},8,["modelValue"]),e($e,{modelValue:t.value.remark,"onUpdate:modelValue":o[11]||(o[11]=l=>t.value.remark=l),modelModifiers:{trim:!0},placeholder:"请输入备注",clearable:""},{prepend:a(()=>[Je]),_:1},8,["modelValue"])]),_:1},512)])]),_:1},512)}}});const Xe=De(Ze,[["__scopeId","data-v-184fa553"]]),ea={class:"text-primary font-weight-black"},aa=u("br",null,null,-1),ka=Y({__name:"speakerManage",setup(n){const{loadDictTree:C,getLabels:y}=j();C(["speak_age_group","speak_gender","speak_provider","speak_lang"]);const h=x({title:"发声人管理"}),U=x([{text:"声音合成",disabled:!1,href:"#"},{text:"发声人管理",disabled:!0,href:"#"}]),_=F({speakName:"",provider:null,lang:null}),f=F({list:[],total:0}),V=x(),t=x(),m=x(),v=F({id:"",name:""}),I=s=>{let p=[];return p.push({text:"删除",color:"error",click(){R(s)}}),p.push({text:"编辑",color:"info",click(){o(s)}}),p},i=async(s={})=>{const[p,l]=await A.get({url:"/voice/speak",showLoading:t.value.el,data:{..._,...s}});l?(f.list=l.list||[],f.total=l.total):(f.list=[],f.total=0)},c=()=>{t.value.query({page:1})},R=s=>{v.name=s.speakCname,v.id=s.id,m.value.show({width:"400px",confirmText:v.name})},T=async(s={})=>{const[p,l]=await A.delete({...s,showSuccess:!0,url:`/voice/speak/${v.id}`});l&&(m.value.hide(),i())},d=()=>{V.value.show({title:"添加发声人",operateType:"add"})},o=s=>{V.value.show({title:"编辑发声人",infos:s,operateType:"edit"})};return se(()=>{i()}),(s,p)=>{const l=S("Select"),J=S("ButtonsInForm"),g=S("el-table-column"),Z=S("router-link"),X=S("ButtonsInTable"),ee=S("TableWithPager"),ae=re("copy");return P(),q(L,null,[e(ie,{title:h.value.title,breadcrumbs:U.value},null,8,["title","breadcrumbs"]),e(pe,null,{default:a(()=>[e(Ne,null,{default:a(()=>[e(D,{cols:"12",lg:"3",md:"4",sm:"6"},{default:a(()=>[e(M,{modelValue:_.speakName,"onUpdate:modelValue":p[0]||(p[0]=r=>_.speakName=r),label:"请输入标识","hide-details":"",clearable:"",onKeyup:ne(c,["enter"]),"onClick:clear":c},null,8,["modelValue","onKeyup"])]),_:1}),e(D,{cols:"12",lg:"3",md:"4",sm:"6"},{default:a(()=>[e(l,{modelValue:_.provider,"onUpdate:modelValue":p[1]||(p[1]=r=>_.provider=r),mapDictionary:{code:"speak_provider"},label:"请选择供应商","hide-details":"",onChange:c},null,8,["modelValue"])]),_:1}),e(D,{cols:"12",lg:"3",md:"4",sm:"6"},{default:a(()=>[e(l,{modelValue:_.lang,"onUpdate:modelValue":p[2]||(p[2]=r=>_.lang=r),mapDictionary:{code:"speak_lang"},label:"请选择语言","hide-details":"",onChange:c},null,8,["modelValue"])]),_:1}),e(D,{cols:"12",lg:"3",md:"4",sm:"6"},{default:a(()=>[e(J,null,{default:a(()=>[e(Ue,{color:"primary",onClick:d},{default:a(()=>[b("添加发声人")]),_:1})]),_:1})]),_:1}),e(D,{cols:"12"},{default:a(()=>[e(me,null,{default:a(()=>[b("修改之后将实时生效，请谨慎操作！")]),_:1})]),_:1}),e(D,{cols:"12"},{default:a(()=>[e(ee,{onQuery:i,ref_key:"tableWithPagerRef",ref:t,infos:f},{default:a(()=>[e(g,{label:"标识",width:"150px","show-overflow-tooltip":""},{default:a(({row:r})=>[de((P(),q("span",null,[b(w(r.speakName),1)])),[[ae,r.speakName]])]),_:1}),e(g,{label:"姓名",prop:"speakCname",width:"100px"}),e(g,{label:"供应",width:"100px"},{default:a(({row:r})=>[u("span",null,w($(y)([["speak_provider",r.provider]])),1)]),_:1}),e(g,{label:"语言",width:"160px"},{default:a(({row:r})=>[u("span",null,w($(y)([["speak_lang",r.lang]])),1)]),_:1}),e(g,{label:"音色",width:"120px"},{default:a(({row:r})=>[u("div",null,w($(y)([["speak_age_group",r.ageGroup],["speak_gender",r.gender]],O=>O.length?O.join("")+"声":"未知")),1)]),_:1}),e(g,{label:"试听","min-width":"330px"},{default:a(({row:r})=>[e(fe,{src:r==null?void 0:r.speakDemo},null,8,["src"])]),_:1}),e(g,{label:"部署",width:"100px"},{default:a(({row:r})=>[e(Z,{to:{path:"/voice-print/synthesis/synthesis-voice",query:{provider:r.provider,lang:r.lang,speakName:r.speakName}},class:"text-info"},{default:a(()=>[b("合成")]),_:2},1032,["to"])]),_:1}),e(g,{label:"备注",prop:"remark","min-width":"200px"}),e(g,{label:"更新时间","min-width":"160px"},{default:a(({row:r})=>[b(w($(ue).dateFormat(r.updatedAt,"YYYY-MM-DD HH:mm:ss")),1)]),_:1}),e(g,{label:"操作",width:"120px",fixed:"right"},{default:a(({row:r})=>[e(X,{buttons:I(r)},null,8,["buttons"])]),_:1})]),_:1},8,["infos"])]),_:1})]),_:1})]),_:1}),e(ce,{ref_key:"refConfirmDelete",ref:m,onSubmit:T},{text:a(()=>[b(" 您将要删除"),u("span",ea,w(v.name),1),b("发声人，删除之后该声音将无法继续合成新的声音。"),aa,b(" 确定要继续吗？ ")]),_:1},512),e(Xe,{ref_key:"createSpeakerPaneRef",ref:V,onSubmit:c},null,512)],64)}}});export{ka as default};