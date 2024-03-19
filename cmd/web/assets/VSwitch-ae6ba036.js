import{p as T,ak as $,ai as p,g as K,A as m,bc as M,aA as q,ab as E,b as G,an as H,ae as g,aj as V,aK as J,af as O,bd as Q,ad as W}from"./index-e2ac1ad9.js";import{r as X,c as h,j as a,L as b,F as Y}from"./utils-41654a3b.js";const Z=T({indeterminate:Boolean,inset:Boolean,flat:Boolean,loading:{type:[Boolean,String],default:!1},...$(),...p()},"VSwitch"),te=K()({name:"VSwitch",inheritAttrs:!1,props:Z(),emits:{"update:focused":e=>!0,"update:modelValue":e=>!0,"update:indeterminate":e=>!0},setup(e,k){let{attrs:C,slots:n}=k;const o=m(e,"indeterminate"),s=m(e,"modelValue"),{loaderClasses:w}=M(e),{isFocused:y,focus:S,blur:P}=q(e),f=X(),A=h(()=>typeof e.loading=="string"&&e.loading!==""?e.loading:e.color),F=E(),_=h(()=>e.id||`switch-${F}`);function x(){o.value&&(o.value=!1)}function B(i){var u,d;i.stopPropagation(),i.preventDefault(),(d=(u=f.value)==null?void 0:u.input)==null||d.click()}return G(()=>{const[i,u]=H(C),d=g.filterProps(e),I=V.filterProps(e);return a(g,b({class:["v-switch",{"v-switch--inset":e.inset},{"v-switch--indeterminate":o.value},w.value,e.class]},i,d,{modelValue:s.value,"onUpdate:modelValue":r=>s.value=r,id:_.value,focused:y.value,style:e.style}),{...n,default:r=>{let{id:L,messagesId:U,isDisabled:j,isReadonly:z,isValid:D}=r;return a(V,b({ref:f},I,{modelValue:s.value,"onUpdate:modelValue":[l=>s.value=l,x],id:L.value,"aria-describedby":U.value,type:"checkbox","aria-checked":o.value?"mixed":void 0,disabled:j.value,readonly:z.value,onFocus:S,onBlur:P},u),{...n,default:l=>{let{backgroundColorClasses:c,backgroundColorStyles:t}=l;return a("div",{class:["v-switch__track",...c.value],style:t.value,onClick:B},null)},input:l=>{let{inputNode:c,icon:t,backgroundColorClasses:N,backgroundColorStyles:R}=l;return a(Y,null,[c,a("div",{class:["v-switch__thumb",{"v-switch__thumb--filled":t||e.loading},e.inset?void 0:N.value],style:e.inset?void 0:R.value},[a(J,null,{default:()=>[e.loading?a(Q,{name:"v-switch",active:!0,color:D.value===!1?void 0:A.value},{default:v=>n.loader?n.loader(v):a(W,{active:v.isActive,color:v.color,indeterminate:!0,size:"16",width:"2"},null)}):t&&a(O,{key:t,icon:t,size:"x-small"},null)]})])])}})}})}),{}}});export{te as V};