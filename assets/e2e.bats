#!/usr/bin/env bats

@test "reject because PodName is blacklisted" {
  run kwctl run annotated-policy.wasm -r assets/test_data/pod.json --settings-json '{"namespace":"default","unsafe_names":["ocatopic","insecure-"],"safe_names":["ocatopic"]}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request rejected
  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*false') -ne 0 ]
  [ $(expr "$output" : ".*The '-pod' name is on the deny list.*") -ne 0 ]
}

@test "accept because the podName is whitelisted" {
  run kwctl run annotated-policy.wasm -r assets/test_data/pod.json --settings-json '{"namespace":"default","unsafe_names":["notsafe-","insecure-"],"safe_names":["ocatopic"]}'
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
}

@test "accept because the whitelist is empty and is not included into the blacklist" {
  run kwctl run annotated-policy.wasm -r assets/test_data/pod.json '{"namespace":"default","unsafe_names":["notsafe-","insecure-"]}'
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
}
