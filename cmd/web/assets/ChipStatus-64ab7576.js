import{c as l,e as r}from"./index-24ceeda1.js";import{x as n,z as i,y as p,k as a,l as c,m,O as u,_ as f,$ as h}from"./utils-15090c58.js";var C=l("circle","IconCircle",[["path",{d:"M12 12m-9 0a9 9 0 1 0 18 0a9 9 0 1 0 -18 0",key:"svg-0"}]]),d=l("clock-hour-4","IconClockHour4",[["path",{d:"M12 12m-9 0a9 9 0 1 0 18 0a9 9 0 1 0 -18 0",key:"svg-0"}],["path",{d:"M12 12l3 2",key:"svg-1"}],["path",{d:"M12 7v5",key:"svg-2"}]]);const k={__name:"ChipStatus",props:{modelValue:{type:Boolean,default(){return!1}}},setup(t){const e=n({label:"",color:"",icon:""});i(e);const s=t;return p(()=>{let{modelValue:o}=s;o?(e.label="启用",e.color="success",e.icon=C):(e.label="停用",e.color="default",e.icon=d)}),(o,v)=>(a(),c(r,{color:e.color,label:"",size:"small"},{default:m(()=>[(a(),c(u(e.icon),{size:"15",class:"mr-1"})),f(h(e.label),1)]),_:1},8,["color"]))}},b=k;export{b as C};
