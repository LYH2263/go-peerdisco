async function load(){
  const m=await fetch('/api/members').then(r=>r.json());
  const s=await fetch('/api/suspicion').then(r=>r.json());
  document.getElementById('members').textContent=JSON.stringify(m,null,2);
  document.getElementById('suspicion').textContent=JSON.stringify(s,null,2);
}
document.getElementById('refresh').onclick=load;
document.getElementById('doJoin').onclick=async()=>{
  const id=document.getElementById('jid').value;
  const addr=document.getElementById('jaddr').value;
  await fetch('/api/join',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({ID:id,Addr:addr,Tags:['ui']})});
  load();
};
load();
setInterval(load,3000);
