import{e as t}from"./index-24ceeda1.js";import{x as s,z as c,y as r,k as n,l as p,m as i,_ as f,$ as m}from"./utils-15090c58.js";const u={__name:"ChipBoolean",props:{modelValue:{type:Boolean,default(){return!1}}},setup(l){const e=s({label:"",color:""});c(e);const a=l;return r(()=>{let{modelValue:o}=a;o?(e.label="是",e.color="success"):(e.label="否",e.color="default")}),(o,_)=>(n(),p(t,{color:e.color,label:"",size:"small"},{default:i(()=>[f(m(e.label),1)]),_:1},8,["color"]))}},b=u;export{b as C};