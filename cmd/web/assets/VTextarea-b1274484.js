import{c as y,r as v,s as p,y as ee,o as te,w as x,K as ae,P as R,j as l,L as S,F as b,Q as _,S as ne,aC as le}from"./utils-15090c58.js";import{p as oe,ae as ue,aF as ie,g as re,b8 as se,A as ce,at as de,b as fe,aH as ve,k as xe,ah as me,ai as G,aI as ge,aJ as he,aK as we,aL as Ve,P as ye}from"./index-24ceeda1.js";const Fe=oe({autoGrow:Boolean,autofocus:Boolean,counter:[Boolean,Number,String],counterValue:Function,prefix:String,placeholder:String,persistentPlaceholder:Boolean,persistentCounter:Boolean,noResize:Boolean,rows:{type:[Number,String],default:5,validator:e=>!isNaN(parseFloat(e))},maxRows:{type:[Number,String],validator:e=>!isNaN(parseFloat(e))},suffix:String,modelModifiers:Object,...ue(),...ie()},"VTextarea"),ke=re()({name:"VTextarea",directives:{Intersect:se},inheritAttrs:!1,props:Fe(),emits:{"click:control":e=>!0,"mousedown:control":e=>!0,"update:focused":e=>!0,"update:modelValue":e=>!0},setup(e,D){let{attrs:F,emit:H,slots:i}=D;const o=ce(e,"modelValue"),{isFocused:f,focus:E,blur:U}=de(e),O=y(()=>typeof e.counterValue=="function"?e.counterValue(o.value):(o.value||"").toString().length),$=y(()=>{if(F.maxlength)return F.maxlength;if(!(!e.counter||typeof e.counter!="number"&&typeof e.counter!="string"))return e.counter});function j(t,n){var a,u;!e.autofocus||!t||(u=(a=n[0].target)==null?void 0:a.focus)==null||u.call(a)}const M=v(),m=v(),A=p(""),g=v(),K=y(()=>e.persistentPlaceholder||f.value||e.active);function P(){var t;g.value!==document.activeElement&&((t=g.value)==null||t.focus()),f.value||E()}function L(t){P(),H("click:control",t)}function J(t){H("mousedown:control",t)}function Q(t){t.stopPropagation(),P(),R(()=>{o.value="",Ve(e["onClick:clear"],t)})}function q(t){var a;const n=t.target;if(o.value=n.value,(a=e.modelModifiers)!=null&&a.trim){const u=[n.selectionStart,n.selectionEnd];R(()=>{n.selectionStart=u[0],n.selectionEnd=u[1]})}}const c=v(),h=v(+e.rows),C=y(()=>["plain","underlined"].includes(e.variant));ee(()=>{e.autoGrow||(h.value=+e.rows)});function d(){e.autoGrow&&R(()=>{if(!c.value||!m.value)return;const t=getComputedStyle(c.value),n=getComputedStyle(m.value.$el),a=parseFloat(t.getPropertyValue("--v-field-padding-top"))+parseFloat(t.getPropertyValue("--v-input-padding-top"))+parseFloat(t.getPropertyValue("--v-field-padding-bottom")),u=c.value.scrollHeight,w=parseFloat(t.lineHeight),k=Math.max(parseFloat(e.rows)*w+a,parseFloat(n.getPropertyValue("--v-input-control-height"))),I=parseFloat(e.maxRows)*w+a||1/0,s=ye(u??0,k,I);h.value=Math.floor((s-a)/w),A.value=xe(s)})}te(d),x(o,d),x(()=>e.rows,d),x(()=>e.maxRows,d),x(()=>e.density,d);let r;return x(c,t=>{t?(r=new ResizeObserver(d),r.observe(c.value)):r==null||r.disconnect()}),ae(()=>{r==null||r.disconnect()}),fe(()=>{const t=!!(i.counter||e.counter||e.counterValue),n=!!(t||i.details),[a,u]=me(F),{modelValue:w,...k}=G.filterProps(e),I=ge(e);return l(G,S({ref:M,modelValue:o.value,"onUpdate:modelValue":s=>o.value=s,class:["v-textarea v-text-field",{"v-textarea--prefixed":e.prefix,"v-textarea--suffixed":e.suffix,"v-text-field--prefixed":e.prefix,"v-text-field--suffixed":e.suffix,"v-textarea--auto-grow":e.autoGrow,"v-textarea--no-resize":e.noResize||e.autoGrow,"v-input--plain-underlined":C.value},e.class],style:e.style},a,k,{centerAffix:h.value===1&&!C.value,focused:f.value}),{...i,default:s=>{let{id:V,isDisabled:B,isDirty:N,isReadonly:W,isValid:X}=s;return l(he,S({ref:m,style:{"--v-textarea-control-height":A.value},onClick:L,onMousedown:J,"onClick:clear":Q,"onClick:prependInner":e["onClick:prependInner"],"onClick:appendInner":e["onClick:appendInner"]},I,{id:V.value,active:K.value||N.value,centerAffix:h.value===1&&!C.value,dirty:N.value||e.dirty,disabled:B.value,focused:f.value,error:X.value===!1}),{...i,default:Y=>{let{props:{class:T,...z}}=Y;return l(b,null,[e.prefix&&l("span",{class:"v-text-field__prefix"},[e.prefix]),_(l("textarea",S({ref:g,class:T,value:o.value,onInput:q,autofocus:e.autofocus,readonly:W.value,disabled:B.value,placeholder:e.placeholder,rows:e.rows,name:e.name,onFocus:P,onBlur:U},z,u),null),[[ne("intersect"),{handler:j},null,{once:!0}]]),e.autoGrow&&_(l("textarea",{class:[T,"v-textarea__sizer"],id:`${z.id}-sizer`,"onUpdate:modelValue":Z=>o.value=Z,ref:c,readonly:!0,"aria-hidden":"true"},null),[[le,o.value]]),e.suffix&&l("span",{class:"v-text-field__suffix"},[e.suffix])])}})},details:n?s=>{var V;return l(b,null,[(V=i.details)==null?void 0:V.call(i,s),t&&l(b,null,[l("span",null,null),l(we,{active:e.persistentCounter||f.value,value:O.value,max:$.value},i.counter)])])}:void 0})}),ve({},M,m,g)}});export{ke as V};