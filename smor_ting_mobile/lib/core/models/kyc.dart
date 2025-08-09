class KycRequest {
  final String country;
  final String idType;
  final String idNumber;
  final String firstName;
  final String lastName;
  final String phone;

  KycRequest({
    required this.country,
    required this.idType,
    required this.idNumber,
    required this.firstName,
    required this.lastName,
    required this.phone,
  });

  Map<String, dynamic> toJson() => {
        'country': country,
        'id_type': idType,
        'id_number': idNumber,
        'first_name': firstName,
        'last_name': lastName,
        'phone': phone,
      };
}

class KycResponse {
  final String status;
  final String reference;

  KycResponse({required this.status, required this.reference});

  factory KycResponse.fromJson(Map<String, dynamic> json) => KycResponse(
        status: json['status'] as String? ?? 'UNKNOWN',
        reference: json['reference'] as String? ?? '',
      );
}


