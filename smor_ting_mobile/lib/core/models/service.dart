import 'package:json_annotation/json_annotation.dart';

part 'service.g.dart';

@JsonSerializable()
class ServiceCategory {
  final String id;
  final String name;
  final String description;
  final String icon;
  final String color;
  @JsonKey(name: 'is_active')
  final bool isActive;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;

  const ServiceCategory({
    required this.id,
    required this.name,
    required this.description,
    required this.icon,
    required this.color,
    required this.isActive,
    required this.createdAt,
    required this.updatedAt,
  });

  factory ServiceCategory.fromJson(Map<String, dynamic> json) => _$ServiceCategoryFromJson(json);
  Map<String, dynamic> toJson() => _$ServiceCategoryToJson(this);
}

@JsonSerializable()
class Service {
  final String id;
  final String name;
  final String description;
  @JsonKey(name: 'category_id')
  final String categoryId;
  @JsonKey(name: 'provider_id')
  final String providerId;
  final double price;
  final String currency;
  final int duration; // in minutes
  final List<String> images;
  @JsonKey(name: 'is_active')
  final bool isActive;
  final double rating;
  @JsonKey(name: 'review_count')
  final int reviewCount;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;

  const Service({
    required this.id,
    required this.name,
    required this.description,
    required this.categoryId,
    required this.providerId,
    required this.price,
    required this.currency,
    required this.duration,
    required this.images,
    required this.isActive,
    required this.rating,
    required this.reviewCount,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Service.fromJson(Map<String, dynamic> json) => _$ServiceFromJson(json);
  Map<String, dynamic> toJson() => _$ServiceToJson(this);

  String get formattedPrice => '$currency ${price.toStringAsFixed(2)}';
  String get formattedDuration {
    if (duration < 60) {
      return '${duration}min';
    } else {
      final hours = duration ~/ 60;
      final minutes = duration % 60;
      return minutes > 0 ? '${hours}h ${minutes}min' : '${hours}h';
    }
  }
}

@JsonSerializable()
class ServiceProvider {
  final String id;
  @JsonKey(name: 'user_id')
  final String userId;
  @JsonKey(name: 'business_name')
  final String businessName;
  final String description;
  final int experience; // years
  final List<String> certifications;
  @JsonKey(name: 'service_areas')
  final List<String> serviceAreas;
  @JsonKey(name: 'is_verified')
  final bool isVerified;
  final double rating;
  @JsonKey(name: 'review_count')
  final int reviewCount;
  @JsonKey(name: 'completed_jobs')
  final int completedJobs;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;

  const ServiceProvider({
    required this.id,
    required this.userId,
    required this.businessName,
    required this.description,
    required this.experience,
    required this.certifications,
    required this.serviceAreas,
    required this.isVerified,
    required this.rating,
    required this.reviewCount,
    required this.completedJobs,
    required this.createdAt,
    required this.updatedAt,
  });

  factory ServiceProvider.fromJson(Map<String, dynamic> json) => _$ServiceProviderFromJson(json);
  Map<String, dynamic> toJson() => _$ServiceProviderToJson(this);

  String get experienceText => '$experience year${experience != 1 ? 's' : ''} experience';
}

enum BookingStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('confirmed')
  confirmed,
  @JsonValue('in_progress')
  inProgress,
  @JsonValue('completed')
  completed,
  @JsonValue('cancelled')
  cancelled,
}

@JsonSerializable()
class Booking {
  final String id;
  @JsonKey(name: 'customer_id')
  final String customerId;
  @JsonKey(name: 'provider_id')
  final String providerId;
  @JsonKey(name: 'service_id')
  final String serviceId;
  final BookingStatus status;
  @JsonKey(name: 'scheduled_date')
  final DateTime scheduledDate;
  @JsonKey(name: 'completed_date')
  final DateTime? completedDate;
  final Address address;
  final String notes;
  @JsonKey(name: 'total_amount')
  final double totalAmount;
  final String currency;
  @JsonKey(name: 'payment_status')
  final String paymentStatus;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;

  const Booking({
    required this.id,
    required this.customerId,
    required this.providerId,
    required this.serviceId,
    required this.status,
    required this.scheduledDate,
    this.completedDate,
    required this.address,
    required this.notes,
    required this.totalAmount,
    required this.currency,
    required this.paymentStatus,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Booking.fromJson(Map<String, dynamic> json) => _$BookingFromJson(json);
  Map<String, dynamic> toJson() => _$BookingToJson(this);

  String get formattedAmount => '$currency ${totalAmount.toStringAsFixed(2)}';
}

@JsonSerializable()
class Address {
  final String street;
  final String city;
  final String county;
  final String country;
  final double latitude;
  final double longitude;

  const Address({
    required this.street,
    required this.city,
    required this.county,
    required this.country,
    required this.latitude,
    required this.longitude,
  });

  factory Address.fromJson(Map<String, dynamic> json) => _$AddressFromJson(json);
  Map<String, dynamic> toJson() => _$AddressToJson(this);

  String get fullAddress => '$street, $city, $county, $country';
}

@JsonSerializable()
class Review {
  final String id;
  @JsonKey(name: 'booking_id')
  final String bookingId;
  @JsonKey(name: 'customer_id')
  final String customerId;
  @JsonKey(name: 'provider_id')
  final String providerId;
  @JsonKey(name: 'service_id')
  final String serviceId;
  final int rating; // 1-5
  final String comment;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;

  const Review({
    required this.id,
    required this.bookingId,
    required this.customerId,
    required this.providerId,
    required this.serviceId,
    required this.rating,
    required this.comment,
    required this.createdAt,
  });

  factory Review.fromJson(Map<String, dynamic> json) => _$ReviewFromJson(json);
  Map<String, dynamic> toJson() => _$ReviewToJson(this);
}