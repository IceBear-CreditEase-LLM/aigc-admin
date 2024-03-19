import{_ as E}from"./NavBack.vue_vue_type_script_setup_true_lang-db3181aa.js";import{x as D,z as I,i as V,c as U,U as O,k as b,l as w,m as e,j as t,_,$ as d,n as f,aB as g,W as o,av as z,al as P,am as R,r as F,y as C,Z as N,O as W,ab as q,au as G,D as J,F as K}from"./utils-766f017d.js";import{_ as M}from"./UiParentCard.vue_vue_type_script_setup_true_lang-3544285d.js";import{C as Q}from"./ChipBoolean-857cecb3.js";import{_ as T,ae as v,ah as $,a0 as Y,ac as X}from"./index-88a5d611.js";import{V as u}from"./VCol-07cdd7be.js";import{V as j}from"./VRow-5c8abc48.js";import{_ as Z}from"./UiChildCard.vue_vue_type_script_setup_true_lang-ee66a36a.js";import{t as L}from"./map-f2c19dca.js";import{r as A,i as B,Z as H}from"./index-8cd50de6.js";import{T as ee}from"./TextLog-f798df3c.js";import{V as te,a as S,b as ae,c as k}from"./VWindowItem-f22f5444.js";import"./IconCircleCheckFilled-a48483eb.js";const h=x=>(P("data-v-f080d916"),x=x(),R(),x),le=h(()=>o("label",null,"创建人",-1)),se=h(()=>o("label",null,"基础模型",-1)),oe=h(()=>o("label",null,"训练模型",-1)),ne=h(()=>o("label",null,"是否Lora训练",-1)),ie=h(()=>o("label",null,"使用GPU",-1)),re=h(()=>o("label",null,"创建时间",-1)),de=h(()=>o("label",null,"训练轮次",-1)),ue=h(()=>o("label",null,"模型后缀",-1)),_e=h(()=>o("label",null,"训练批次",-1)),ce=h(()=>o("label",null,"学习率",-1)),pe=h(()=>o("label",null,"开始时间",-1)),me=h(()=>o("label",null,"备注",-1)),fe=h(()=>o("label",null,"文件",-1)),ve=h(()=>o("label",null,"完成时间",-1)),he=h(()=>o("label",null,"模型最大长度",-1)),xe={__name:"FineTuningBaseInfo",setup(x){const c=D({style:{},formData:{}});I(c);const s=V("provideFineTuningDetail"),p=n=>{z.downloadByUrl({fileUrl:n.fileUrl,suffixName:"jsonl"})},l=U(()=>s.rawData);return(n,a)=>{const i=O("router-link");return b(),w(j,{class:"my-form waterfall"},{default:e(()=>[t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[le]),default:e(()=>[_(" "+d(l.value.trainPublisher),1)]),_:1})]),_:1}),t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[se]),default:e(()=>[t(i,{to:"",class:"link"},{default:e(()=>[_(d(l.value.baseModel),1)]),_:1})]),_:1})]),_:1}),t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[oe]),default:e(()=>[_(" "+d(l.value.fineTunedModel),1)]),_:1})]),_:1}),t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[ne]),default:e(()=>[t(Q,{modelValue:l.value.lora,"onUpdate:modelValue":a[0]||(a[0]=r=>l.value.lora=r)},null,8,["modelValue"])]),_:1})]),_:1}),t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[ie]),default:e(()=>[_(" "+d(l.value.procPerNode),1)]),_:1})]),_:1}),t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[re]),default:e(()=>[_(" "+d(f(g).dateFormat(l.value.createdAt,"YYYY-MM-DD HH:mm:ss")),1)]),_:1})]),_:1}),t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[de]),default:e(()=>[_(" "+d(l.value.trainEpoch),1)]),_:1})]),_:1}),t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[ue]),default:e(()=>[_(" "+d(l.value.suffix),1)]),_:1})]),_:1}),t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[_e]),default:e(()=>[_(" "+d(l.value.trainBatchSize)+"次训练,"+d(l.value.evalBatchSize)+"次评估,"+d(l.value.accumulationSteps)+" 次梯度累加 ",1)]),_:1})]),_:1}),t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[ce]),default:e(()=>[_(" "+d(f(g).toScientfic(l.value.learningRate)),1)]),_:1})]),_:1}),t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[pe]),default:e(()=>[_(" "+d(f(g).dateFormat(l.value.startTrainTime,"YYYY-MM-DD HH:mm:ss")),1)]),_:1})]),_:1}),t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[me]),default:e(()=>[_(" "+d(l.value.remark),1)]),_:1})]),_:1}),t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[fe]),default:e(()=>[o("a",{class:"link1 line1",onClick:a[1]||(a[1]=r=>p(l.value))},d(l.value.fileId),1)]),_:1})]),_:1}),t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[ve]),default:e(()=>[_(" "+d(f(g).dateFormat(l.value.finishedAt,"YYYY-MM-DD HH:mm:ss")),1)]),_:1})]),_:1}),t(u,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[t(v,{"hide-details":""},{prepend:e(()=>[he]),default:e(()=>[_(" "+d(l.value.modelMaxLength),1)]),_:1})]),_:1})]),_:1})}}},ge=T(xe,[["__scopeId","data-v-f080d916"]]);const be={__name:"Chart1",setup(x){D({style:{},formData:{}});const c=F();let s;const p=V("provideFineTuningDetail"),l=()=>{var r;let{epoch:n}=((r=p.rawData)==null?void 0:r.trainAnalysis)||{},a=[],i=[];n.list.forEach((m,y)=>{a.push(y),i.push(m.value)}),A("ZH",H),s=B(c.value,null,{locale:"ZH"}),s.setOption({title:{text:"train/epoch",x:"center"},toolbox:{feature:{saveAsImage:{},dataView:{},dataZoom:{}}},tooltip:{trigger:"axis",axisPointer:{type:"cross"}},xAxis:{type:"category",data:a,axisLine:{onZero:!1},axisLabel:{formatter(m){return m==0?"":m},interval:Math.floor((a.length-1)/10)}},yAxis:{type:"value"},series:[{data:i,type:"line",smooth:!0}]})};return C(()=>{var a;let{epoch:n}=((a=p.rawData)==null?void 0:a.trainAnalysis)||{};n&&l()}),$(c,n=>{s==null||s.resize()}),(n,a)=>(b(),w(Y,null,{default:e(()=>[o("div",{ref_key:"refBox",ref:c,class:"chart-item w-100"},null,512)]),_:1}))}},De=T(be,[["__scopeId","data-v-a3a8208b"]]);const ye={__name:"Chart2",setup(x){D({style:{},formData:{}});const c=F();let s;const p=V("provideFineTuningDetail"),l=()=>{var r;let{loss:n}=((r=p.rawData)==null?void 0:r.trainAnalysis)||{},a=[],i=[];n.list.forEach((m,y)=>{a.push(y),i.push(m.value)}),A("ZH",H),s=B(c.value,null,{locale:"ZH"}),s.setOption({title:{text:"train/loss",x:"center"},toolbox:{feature:{saveAsImage:{},dataView:{},dataZoom:{}}},tooltip:{trigger:"axis",axisPointer:{type:"cross"}},xAxis:{type:"category",data:a,axisLine:{onZero:!1},axisLabel:{formatter(m){return m==0?"":m},interval:Math.floor((a.length-1)/10)}},yAxis:{type:"value",axisLine:{onZero:!1}},series:[{data:i,type:"line",smooth:!0}]})};return C(()=>{var a;let{loss:n}=((a=p.rawData)==null?void 0:a.trainAnalysis)||{};n&&l()}),$(c,n=>{s==null||s.resize()}),(n,a)=>(b(),w(Y,null,{default:e(()=>[o("div",{ref_key:"refBox",ref:c,class:"chart-item w-100"},null,512)]),_:1}))}},Ve=T(ye,[["__scopeId","data-v-59df7dc9"]]);const we={__name:"Chart3",setup(x){D({style:{},formData:{}});const c=F();let s;const p=V("provideFineTuningDetail"),l=()=>{var r;let{learningRate:n}=((r=p.rawData)==null?void 0:r.trainAnalysis)||{},a=[],i=[];n.list.forEach((m,y)=>{a.push(y),i.push(m.value)}),A("ZH",H),s=B(c.value,null,{locale:"ZH"}),s.setOption({title:{text:"train/learningRate",x:"center"},toolbox:{feature:{saveAsImage:{},dataView:{},dataZoom:{}}},tooltip:{trigger:"axis",axisPointer:{type:"cross"}},xAxis:{type:"category",data:a,axisLine:{onZero:!1},axisLabel:{formatter(m){return m==0?"":m},interval:Math.floor((a.length-1)/10)}},yAxis:{type:"value"},series:[{data:i,type:"line",smooth:!0}]})};return C(()=>{var a;let{learningRate:n}=((a=p.rawData)==null?void 0:a.trainAnalysis)||{};n&&l()}),$(c,n=>{s==null||s.resize()}),(n,a)=>(b(),w(Y,null,{default:e(()=>[o("div",{ref_key:"refBox",ref:c,class:"chart-item w-100"},null,512)]),_:1}))}},Te=T(we,[["__scopeId","data-v-64fef634"]]);const Ie=x=>(P("data-v-b4bddbcb"),x=x(),R(),x),Fe={class:"pb-4"},Ce={style:{padding:"20px 60px 10px"}},$e={class:"d-flex points justify-space-between"},Ye={class:"item text-center"},Ae=Ie(()=>o("div",{class:"text-h6"},"创建训练任务",-1)),Be={class:"text-subtitle-1 text-medium-emphasis mt-1"},He={class:"item text-center"},Me={class:"text-h6 hv-center"},Ze={class:"text-subtitle-1 text-medium-emphasis mt-1"},Le={__name:"TabFineTuningDetail",setup(x){const c=V("provideFineTuningDetail"),s=D({process:{value:0,icon:"",striped:!1,color:"",iconColor:""}}),{process:p}=I(s),l=U(()=>{let n=c.rawData,{process:a}=s;a.value=n.process*100,a.valueCN=g.toPercent(n.process,2);let{trainStatus:i}=n;return["running"].includes(i)&&(a.striped=!0),i=="cancel"?a.color="#ccc":i=="failed"?a.color="rgb(var(--v-theme-error))":a.color="rgb(var(--v-theme-info))",n});return(n,a)=>(b(),N("div",null,[t(Z,{title:"训练进度"},{default:e(()=>{var i,r;return[o("div",Fe,[o("div",Ce,[t(X,{modelValue:f(p).value,"onUpdate:modelValue":a[0]||(a[0]=m=>f(p).value=m),color:f(p).color,height:"25",striped:f(p).striped},{default:e(({value:m})=>[o("strong",null,d(f(p).valueCN),1)]),_:1},8,["modelValue","color","striped"])]),o("div",$e,[o("div",Ye,[Ae,o("div",Be,d(f(g).dateFormat(l.value.createdAt,"YYYY-MM-DD HH:mm:ss")),1)]),o("div",He,[o("div",Me,[(b(),w(W((i=f(L)[l.value.trainStatus])==null?void 0:i.icon),{size:24,color:(r=f(L)[l.value.trainStatus])==null?void 0:r.iconColor,class:"mr-1"},null,8,["color"])),_(" "+d(f(q).localData.local_trainStatus[l.value.trainStatus]),1)]),o("div",Ze,d(f(g).dateFormat(l.value.finishedAt,"YYYY-MM-DD HH:mm:ss")),1)])])])]}),_:1}),t(Z,{title:"训练统计",class:"mt-4"},{default:e(()=>[t(j,null,{default:e(()=>[t(u,{md:"12",sm:"12"},{default:e(()=>[t(De)]),_:1}),t(u,{md:"12",sm:"12"},{default:e(()=>[t(Ve)]),_:1}),t(u,{md:"12",sm:"12"},{default:e(()=>[t(Te)]),_:1})]),_:1})]),_:1})]))}},Se=T(Le,[["__scopeId","data-v-b4bddbcb"]]),Ke={__name:"fineTuningDetail",setup(x){const c=G(),s=D({tabIndex:"",style:{},rawData:{}}),{style:p,rawData:l}=I(s);return J("provideFineTuningDetail",s),(async()=>{let{jobId:a}=c.query,[i,r]=await z.get({showLoading:!0,url:`/api/finetuning/${a}`});r&&(s.rawData=r)})(),(a,i)=>(b(),N(K,null,[t(E,{backUrl:"/model/fine-tuning/list"},{default:e(()=>[_("微调详情")]),_:1}),t(M,{class:"mt-4"},{header:e(()=>[_(" 任务ID："+d(f(l).jobId),1)]),default:e(()=>[t(ge)]),_:1}),t(M,{class:"mt-5"},{header:e(()=>[t(te,{modelValue:s.tabIndex,"onUpdate:modelValue":i[0]||(i[0]=r=>s.tabIndex=r),"align-tabs":"start",color:"primary"},{default:e(()=>[t(S,{value:1},{default:e(()=>[_("详情")]),_:1}),t(S,{value:2},{default:e(()=>[_("日志")]),_:1})]),_:1},8,["modelValue"])]),default:e(()=>[t(ae,{modelValue:s.tabIndex,"onUpdate:modelValue":i[2]||(i[2]=r=>s.tabIndex=r)},{default:e(()=>[t(k,{value:1},{default:e(()=>[t(Se)]),_:1}),t(k,{value:2},{default:e(()=>[t(ee,{modelValue:f(l).trainLog,"onUpdate:modelValue":i[1]||(i[1]=r=>f(l).trainLog=r),style:{height:"600px"},idDone:!0},null,8,["modelValue"])]),_:1})]),_:1},8,["modelValue"])]),_:1})],64))}};export{Ke as default};
