find deploy/clusters -name $CUSTOMER -print0 | xargs -0 rm -rf
holos render platform -t flatten -t step3 \
  --selector customer=$CUSTOMER
