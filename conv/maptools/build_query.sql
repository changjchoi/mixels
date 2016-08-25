select 
two, three, coalesce(nullif(five,''), four),
ten,
case when length(five) > 0 then 
  case when twentyseven == '1' then 
    printf('(%s,%s)', four, twentysix) 
  else 
    printf('(%s)', four) 
  end 
else 
  case when twentyseven == '1' then 
    printf('(%s)', twentysix) 
  else 
    '' 
  end 
end,
case when thirteen == '0' then 
  case when eleven == '1' then 
    printf('지하 %s', twelve) 
  else 
    printf('%s', twelve) 
  end 
else 
  case when eleven == '1' then 
    printf('지하 %s-%s', twelve, thirteen) 
  else 
    printf('%s-%s', twelve, thirteen) 
  end 
end  
from build_seoul limit 1;

