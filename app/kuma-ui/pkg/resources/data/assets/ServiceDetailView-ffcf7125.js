import{u as q}from"./vue-router-573afc44.js";import{k as x,u as L}from"./store-2fff246d.js";import{Q as T}from"./QueryParameter-70743f73.js";import{_ as A}from"./EmptyBlock.vue_vue_type_script_setup_true_lang-4047971f.js";import{E as B}from"./ErrorBlock-a792d9a1.js";import{_ as V}from"./LoadingBlock.vue_vue_type_script_setup_true_lang-d3176fee.js";import{D as $}from"./DataPlaneList-b3eb9311.js";import{S as z}from"./ServiceSummary-7700bf27.js";import{d as P,h as F,g as I,e as j,f as R,a as m,b as C,F as Q,o as n,r as u,y as O}from"./runtime-dom.esm-bundler-91b41870.js";import"./vuex.esm-bundler-df5bd11e.js";import"./constants-31fdaf55.js";import"./kongponents.es-3df60cd6.js";import"./_plugin-vue_export-helper-c27b6911.js";import"./datadogLogEvents-4578cfa7.js";import"./ContentWrapper-b0fa1a61.js";import"./DataOverview-66f4d3ae.js";import"./StatusBadge-81464ebd.js";import"./TagList-91d1133a.js";import"./YamlView.vue_vue_type_script_setup_true_lang-6d3c03d3.js";import"./index-a8834e9c.js";import"./CodeBlock.vue_vue_type_style_index_0_lang-525d6c39.js";import"./_commonjsHelpers-87174ba5.js";const J={class:"component-frame"},M=P({__name:"ServiceDetails",props:{service:{type:Object,required:!0},externalService:{type:Object,required:!1,default:null},dataPlaneOverviews:{type:Array,required:!1,default:null},dppFilterFields:{type:Object,required:!0},selectedDppName:{type:String,required:!1,default:null}},emits:["load-dataplane-overviews"],setup(d,{emit:y}){const s=d;function e(f,i){y("load-dataplane-overviews",f,i)}return(f,i)=>{var o;return n(),F(Q,null,[I("div",J,[j(z,{service:s.service,"external-service":d.externalService},null,8,["service","external-service"])]),R(),s.dataPlaneOverviews!==null?(n(),m($,{key:0,class:"mt-4","data-plane-overviews":s.dataPlaneOverviews,"dpp-filter-fields":s.dppFilterFields,"selected-dpp-name":s.selectedDppName,"is-gateway-view":((o=s.dataPlaneOverviews[0])==null?void 0:o.dataplane.networking.gateway)!==void 0,onLoadData:e},null,8,["data-plane-overviews","dpp-filter-fields","selected-dpp-name","is-gateway-view"])):C("",!0)],64)}}}),W={class:"service-details"},fe=P({__name:"ServiceDetailView",props:{selectedDppName:{type:String,required:!1,default:null}},setup(d){const y=d,s={name:{description:"filter by name or parts of a name"},protocol:{description:"filter by “kuma.io/protocol” value"},tag:{description:"filter by tags (e.g. “tag: version:2”)"},zone:{description:"filter by “kuma.io/zone” value"}},e=q(),f=L(),i=u(null),o=u(null),g=u(null),_=u(!0),w=u(null);O(()=>e.params.mesh,function(){e.name==="service-detail-view"&&h(0)}),O(()=>e.params.name,function(){e.name==="service-detail-view"&&h(0)});function N(){f.dispatch("updatePageTitle",e.params.service);const a=T.get("filterFields"),l=a!==null?JSON.parse(a):{};h(0,l)}N();async function h(a,l={}){_.value=!0,w.value=null,i.value=null,o.value=null,g.value=null;const c=e.params.mesh,p=e.params.service;try{i.value=await x.getServiceInsight({mesh:c,name:p}),i.value.serviceType==="external"?o.value=await x.getExternalService({mesh:c,name:p}):await S(a,l)}catch(t){t instanceof Error?w.value=t:console.error(t)}finally{_.value=!1}}async function S(a,l){const c=e.params.mesh,p=e.params.service;try{const t=b(p,a,l),r=await x.getAllDataplaneOverviewsFromMesh({mesh:c},t);g.value=r.items??[]}catch{g.value=null}}function b(a,l,c){const t=`kuma.io/service:${a}`,r={...c,offset:l,size:50};if(r.tag){const D=Array.isArray(r.tag)?r.tag:[r.tag],k=[];for(const[v,E]of D.entries())E.startsWith("kuma.io/service:")&&k.push(v);for(let v=k.length-1;v===0;v--)D.splice(k[v],1);r.tag=D.concat(t)}else r.tag=t;return r}return(a,l)=>(n(),F("div",W,[_.value?(n(),m(V,{key:0})):w.value!==null?(n(),m(B,{key:1,error:w.value},null,8,["error"])):i.value===null?(n(),m(A,{key:2})):(n(),m(M,{key:3,service:i.value,"data-plane-overviews":g.value,"external-service":o.value,"dpp-filter-fields":s,"selected-dpp-name":y.selectedDppName,onLoadDataplaneOverviews:S},null,8,["service","data-plane-overviews","external-service","selected-dpp-name"]))]))}});export{fe as default};
