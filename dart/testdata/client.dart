/// Code generated from jsonrpc schema by rpcgen v2.4.sdsd2; DO NOT EDIT.

import 'package:json_annotation/json_annotation.dart';
import 'package:smd_annotations/annotations.dart';

part 'client.g.dart';

@JsonSerializable()
class Campaign {
  @JsonKey(name: 'groups')
  final List<Group> groups;
  @JsonKey(name: 'id')
  final int id;

  Campaign({
    required this.groups,
    required this.id,
  });

  Map<String, dynamic> toJson() => _$CampaignToJson(this);

  factory Campaign.fromJson(Map<String, dynamic> json) =>
      _$CampaignFromJson(json);
}

@JsonSerializable()
class Group {
  @JsonKey(name: 'child')
  final Group? child;
  @JsonKey(name: 'groups')
  final List<Group> groups;
  @JsonKey(name: 'id')
  final int id;
  @JsonKey(name: 'nodes')
  final List<Group> nodes;
  @JsonKey(name: 'sub')
  final SubGroup sub;
  @JsonKey(name: 'title')
  final String title;

  Group({
    this.child,
    required this.groups,
    required this.id,
    required this.nodes,
    required this.sub,
    required this.title,
  });

  Map<String, dynamic> toJson() => _$GroupToJson(this);

  factory Group.fromJson(Map<String, dynamic> json) =>
      _$GroupFromJson(json);
}

@JsonSerializable()
class SubGroup {
  @JsonKey(name: 'id')
  final int id;
  @JsonKey(name: 'nodes')
  final List<Group> nodes;
  @JsonKey(name: 'title')
  final String title;

  SubGroup({
    required this.id,
    required this.nodes,
    required this.title,
  });

  Map<String, dynamic> toJson() => _$SubGroupToJson(this);

  factory SubGroup.fromJson(Map<String, dynamic> json) =>
      _$SubGroupFromJson(json);
}

@JsonSerializable()
class CatalogueFirstParams {
  @JsonKey(name: 'groups')
  final List<Group> groups;

  CatalogueFirstParams({
    required this.groups,
  });

  Map<String, dynamic> toJson() => _$CatalogueFirstParamsToJson(this);

  factory CatalogueFirstParams.fromJson(Map<String, dynamic> json) =>
      _$CatalogueFirstParamsFromJson(json);
}

@JsonSerializable()
class CatalogueSecondParams {
  @JsonKey(name: 'campaigns')
  final List<Campaign> campaigns;

  CatalogueSecondParams({
    required this.campaigns,
  });

  Map<String, dynamic> toJson() => _$CatalogueSecondParamsToJson(this);

  factory CatalogueSecondParams.fromJson(Map<String, dynamic> json) =>
      _$CatalogueSecondParamsFromJson(json);
}



CatalogueRPC catalogueRPCInstance({required RPC rpc}) => _CatalogueRPC(rpc: rpc);
@RPCNamespace('catalogue')
abstract class CatalogueRPC {
  @RPCMethod('First')
  Future<bool> first(CatalogueFirstParams params);
  @RPCMethod('Second')
  Future<bool> second(CatalogueSecondParams params);
  @RPCMethod('Third')
  Future<Campaign> third();
}
