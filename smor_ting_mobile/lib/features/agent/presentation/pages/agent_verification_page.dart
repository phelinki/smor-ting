import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

enum VerificationStep {
  personalInfo,
  idUpload,
  certifications,
  backgroundCheck,
  review,
}

class AgentVerificationPage extends ConsumerStatefulWidget {
  const AgentVerificationPage({super.key});

  @override
  ConsumerState<AgentVerificationPage> createState() => _AgentVerificationPageState();
}

class _AgentVerificationPageState extends ConsumerState<AgentVerificationPage> {
  VerificationStep currentStep = VerificationStep.personalInfo;
  final _formKey = GlobalKey<FormState>();
  
  // Personal Info
  final _firstNameController = TextEditingController();
  final _lastNameController = TextEditingController();
  final _phoneController = TextEditingController();
  final _addressController = TextEditingController();
  String _selectedGender = 'Male';
  
  // ID Upload
  String? _idFrontImage;
  String? _idBackImage;
  String? _selfieImage;
  
  // Certifications
  List<String> _selectedCertifications = [];
  final _certificationController = TextEditingController();
  
  // Background Check
  bool _hasCriminalRecord = false;
  bool _hasValidLicense = true;
  bool _agreesToTerms = false;

  final List<String> _availableCertifications = [
    'Plumbing License',
    'Electrical License',
    'HVAC Certification',
    'Carpentry License',
    'Cleaning Certification',
    'Security License',
  ];

  @override
  void dispose() {
    _firstNameController.dispose();
    _lastNameController.dispose();
    _phoneController.dispose();
    _addressController.dispose();
    _certificationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Color(0xFF007AFF)),
          onPressed: () => context.pop(),
        ),
        title: const Text(
          'Agent Verification',
          style: TextStyle(
            color: Colors.black,
            fontSize: 18,
            fontWeight: FontWeight.w600,
          ),
        ),
        centerTitle: true,
      ),
      body: Column(
        children: [
          // Progress Bar
          Container(
            padding: const EdgeInsets.all(16),
            child: Column(
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      'Step ${_getCurrentStepNumber()} of 5',
                      style: const TextStyle(
                        fontSize: 14,
                        color: Color(0xFF8E8E93),
                      ),
                    ),
                    Text(
                      '${_getProgressPercentage()}%',
                      style: const TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
                        color: Color(0xFF007AFF),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 8),
                LinearProgressIndicator(
                  value: _getProgressValue(),
                  backgroundColor: const Color(0xFFE5E5EA),
                  valueColor: const AlwaysStoppedAnimation<Color>(Color(0xFF007AFF)),
                ),
                const SizedBox(height: 16),
                Text(
                  _getStepTitle(),
                  style: const TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.w600,
                    color: Colors.black,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  _getStepDescription(),
                  style: const TextStyle(
                    fontSize: 14,
                    color: Color(0xFF8E8E93),
                  ),
                  textAlign: TextAlign.center,
                ),
              ],
            ),
          ),
          
          // Form Content
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(16),
              child: Form(
                key: _formKey,
                child: _buildCurrentStep(),
              ),
            ),
          ),
          
          // Navigation Buttons
          Container(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                if (currentStep != VerificationStep.personalInfo)
                  Expanded(
                    child: OutlinedButton(
                      onPressed: _previousStep,
                      style: OutlinedButton.styleFrom(
                        foregroundColor: const Color(0xFF007AFF),
                        side: const BorderSide(color: Color(0xFF007AFF)),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                        padding: const EdgeInsets.symmetric(vertical: 16),
                      ),
                      child: const Text('Previous'),
                    ),
                  ),
                if (currentStep != VerificationStep.personalInfo)
                  const SizedBox(width: 16),
                Expanded(
                  child: ElevatedButton(
                    onPressed: _nextStep,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFFFF3B30),
                      foregroundColor: Colors.white,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                      padding: const EdgeInsets.symmetric(vertical: 16),
                    ),
                    child: Text(
                      currentStep == VerificationStep.review ? 'Submit' : 'Next',
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCurrentStep() {
    switch (currentStep) {
      case VerificationStep.personalInfo:
        return _buildPersonalInfoStep();
      case VerificationStep.idUpload:
        return _buildIdUploadStep();
      case VerificationStep.certifications:
        return _buildCertificationsStep();
      case VerificationStep.backgroundCheck:
        return _buildBackgroundCheckStep();
      case VerificationStep.review:
        return _buildReviewStep();
    }
  }

  Widget _buildPersonalInfoStep() {
    return Column(
      children: [
        TextFormField(
          controller: _firstNameController,
          decoration: _buildInputDecoration('First Name', Icons.person),
          validator: (value) {
            if (value == null || value.isEmpty) {
              return 'Please enter your first name';
            }
            return null;
          },
        ),
        const SizedBox(height: 16),
        TextFormField(
          controller: _lastNameController,
          decoration: _buildInputDecoration('Last Name', Icons.person),
          validator: (value) {
            if (value == null || value.isEmpty) {
              return 'Please enter your last name';
            }
            return null;
          },
        ),
        const SizedBox(height: 16),
        TextFormField(
          controller: _phoneController,
          keyboardType: TextInputType.phone,
          decoration: _buildInputDecoration('Phone Number', Icons.phone),
          validator: (value) {
            if (value == null || value.isEmpty) {
              return 'Please enter your phone number';
            }
            return null;
          },
        ),
        const SizedBox(height: 16),
        TextFormField(
          controller: _addressController,
          maxLines: 3,
          decoration: _buildInputDecoration('Address', Icons.location_on),
          validator: (value) {
            if (value == null || value.isEmpty) {
              return 'Please enter your address';
            }
            return null;
          },
        ),
        const SizedBox(height: 16),
        DropdownButtonFormField<String>(
          value: _selectedGender,
          decoration: _buildInputDecoration('Gender', Icons.person_outline),
          items: ['Male', 'Female', 'Other'].map((String value) {
            return DropdownMenuItem<String>(
              value: value,
              child: Text(value),
            );
          }).toList(),
          onChanged: (String? newValue) {
            setState(() {
              _selectedGender = newValue!;
            });
          },
        ),
      ],
    );
  }

  Widget _buildIdUploadStep() {
    return Column(
      children: [
        _buildImageUploadCard(
          'ID Card (Front)',
          'Upload the front of your ID card',
          _idFrontImage,
          () => _uploadImage('front'),
        ),
        const SizedBox(height: 16),
        _buildImageUploadCard(
          'ID Card (Back)',
          'Upload the back of your ID card',
          _idBackImage,
          () => _uploadImage('back'),
        ),
        const SizedBox(height: 16),
        _buildImageUploadCard(
          'Selfie',
          'Take a selfie holding your ID',
          _selfieImage,
          () => _uploadImage('selfie'),
        ),
      ],
    );
  }

  Widget _buildCertificationsStep() {
    return Column(
      children: [
        const Text(
          'Select your certifications and licenses',
          style: TextStyle(
            fontSize: 16,
            color: Color(0xFF8E8E93),
          ),
        ),
        const SizedBox(height: 16),
        ...(_availableCertifications.map((cert) => CheckboxListTile(
          title: Text(cert),
          value: _selectedCertifications.contains(cert),
          onChanged: (bool? value) {
            setState(() {
              if (value == true) {
                _selectedCertifications.add(cert);
              } else {
                _selectedCertifications.remove(cert);
              }
            });
          },
          activeColor: const Color(0xFF007AFF),
        ))),
        const SizedBox(height: 16),
        TextFormField(
          controller: _certificationController,
          decoration: _buildInputDecoration('Other Certifications', Icons.add).copyWith(
            hintText: 'Add any other certifications...',
          ),
        ),
      ],
    );
  }

  Widget _buildBackgroundCheckStep() {
    return Column(
      children: [
        SwitchListTile(
          title: const Text('I have a valid professional license'),
          subtitle: const Text('Confirm you have the required licenses'),
          value: _hasValidLicense,
          onChanged: (bool value) {
            setState(() {
              _hasValidLicense = value;
            });
          },
          activeColor: const Color(0xFF007AFF),
        ),
        SwitchListTile(
          title: const Text('I have no criminal record'),
          subtitle: const Text('Confirm you have a clean background'),
          value: !_hasCriminalRecord,
          onChanged: (bool value) {
            setState(() {
              _hasCriminalRecord = !value;
            });
          },
          activeColor: const Color(0xFF007AFF),
        ),
        SwitchListTile(
          title: const Text('I agree to the terms and conditions'),
          subtitle: const Text('Read and accept our terms'),
          value: _agreesToTerms,
          onChanged: (bool value) {
            setState(() {
              _agreesToTerms = value;
            });
          },
          activeColor: const Color(0xFF007AFF),
        ),
        const SizedBox(height: 16),
        Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: const Color(0xFFF2F2F7),
            borderRadius: BorderRadius.circular(12),
          ),
          child: const Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Background Check Notice',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                  color: Colors.black,
                ),
              ),
              SizedBox(height: 8),
              Text(
                'By proceeding, you authorize us to conduct a background check. This helps ensure the safety and trust of our platform.',
                style: TextStyle(
                  fontSize: 14,
                  color: Color(0xFF8E8E93),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildReviewStep() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Review Your Information',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.w600,
            color: Colors.black,
          ),
        ),
        const SizedBox(height: 16),
        _buildReviewSection('Personal Information', [
          'Name: ${_firstNameController.text} ${_lastNameController.text}',
          'Phone: ${_phoneController.text}',
          'Address: ${_addressController.text}',
          'Gender: $_selectedGender',
        ]),
        const SizedBox(height: 16),
        _buildReviewSection('ID Documents', [
          'ID Front: ${_idFrontImage != null ? "Uploaded" : "Not uploaded"}',
          'ID Back: ${_idBackImage != null ? "Uploaded" : "Not uploaded"}',
          'Selfie: ${_selfieImage != null ? "Uploaded" : "Not uploaded"}',
        ]),
        const SizedBox(height: 16),
        _buildReviewSection('Certifications', [
          'Selected: ${_selectedCertifications.join(", ")}',
          'Other: ${_certificationController.text.isNotEmpty ? _certificationController.text : "None"}',
        ]),
        const SizedBox(height: 16),
        _buildReviewSection('Background Check', [
          'Valid License: ${_hasValidLicense ? "Yes" : "No"}',
          'Clean Record: ${!_hasCriminalRecord ? "Yes" : "No"}',
          'Terms Accepted: ${_agreesToTerms ? "Yes" : "No"}',
        ]),
      ],
    );
  }

  Widget _buildReviewSection(String title, List<String> items) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFFF2F2F7),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            title,
            style: const TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: Colors.black,
            ),
          ),
          const SizedBox(height: 8),
          ...(items.map((item) => Padding(
            padding: const EdgeInsets.only(bottom: 4),
            child: Text(
              item,
              style: const TextStyle(
                fontSize: 14,
                color: Color(0xFF8E8E93),
              ),
            ),
          ))),
        ],
      ),
    );
  }

  Widget _buildImageUploadCard(String title, String subtitle, String? image, VoidCallback onTap) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        border: Border.all(color: const Color(0xFFE5E5EA)),
        borderRadius: BorderRadius.circular(12),
      ),
      child: InkWell(
        onTap: onTap,
        child: Column(
          children: [
            if (image != null)
              Container(
                width: 100,
                height: 100,
                decoration: BoxDecoration(
                  color: const Color(0xFF34C759),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: const Icon(
                  Icons.check,
                  color: Colors.white,
                  size: 40,
                ),
              )
            else
              Container(
                width: 100,
                height: 100,
                decoration: BoxDecoration(
                  color: const Color(0xFFF2F2F7),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: const Icon(
                  Icons.camera_alt,
                  color: Color(0xFF8E8E93),
                  size: 40,
                ),
              ),
            const SizedBox(height: 12),
            Text(
              title,
              style: const TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: Colors.black,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              subtitle,
              style: const TextStyle(
                fontSize: 14,
                color: Color(0xFF8E8E93),
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  InputDecoration _buildInputDecoration(String label, IconData icon) {
    return InputDecoration(
      labelText: label,
      prefixIcon: Icon(icon, color: const Color(0xFF007AFF)),
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: Color(0xFFE5E5EA)),
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: Color(0xFFE5E5EA)),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: Color(0xFF007AFF), width: 2),
      ),
      filled: true,
      fillColor: const Color(0xFFF2F2F7),
    );
  }

  void _uploadImage(String type) {
    // Simulate image upload
    setState(() {
      switch (type) {
        case 'front':
          _idFrontImage = 'uploaded';
          break;
        case 'back':
          _idBackImage = 'uploaded';
          break;
        case 'selfie':
          _selfieImage = 'uploaded';
          break;
      }
    });
    _showSnackBar('Image uploaded successfully');
  }

  void _nextStep() {
    if (currentStep == VerificationStep.review) {
      _submitVerification();
    } else {
      if (_formKey.currentState!.validate()) {
        setState(() {
          currentStep = VerificationStep.values[currentStep.index + 1];
        });
      }
    }
  }

  void _previousStep() {
    if (currentStep != VerificationStep.personalInfo) {
      setState(() {
        currentStep = VerificationStep.values[currentStep.index - 1];
      });
    }
  }

  void _submitVerification() {
    // Simulate API call
    _showLoadingDialog();
    
    Future.delayed(const Duration(seconds: 2), () {
      Navigator.of(context).pop();
      _showSuccessDialog();
    });
  }

  void _showLoadingDialog() {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return const AlertDialog(
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              CircularProgressIndicator(),
              SizedBox(height: 16),
              Text('Submitting verification...'),
            ],
          ),
        );
      },
    );
  }

  void _showSuccessDialog() {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return AlertDialog(
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Icon(
                Icons.check_circle,
                color: Color(0xFF34C759),
                size: 64,
              ),
              const SizedBox(height: 16),
              const Text(
                'Verification Submitted!',
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.w600,
                  color: Colors.black,
                ),
              ),
              const SizedBox(height: 8),
              const Text(
                'Your verification is being reviewed. You will be notified once approved.',
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 14,
                  color: Color(0xFF8E8E93),
                ),
              ),
              const SizedBox(height: 24),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: () {
                    Navigator.of(context).pop();
                    context.go('/agent-dashboard');
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFFFF3B30),
                    foregroundColor: Colors.white,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: const Text('Continue'),
                ),
              ),
            ],
          ),
        );
      },
    );
  }

  int _getCurrentStepNumber() {
    return currentStep.index + 1;
  }

  int _getProgressPercentage() {
    return ((currentStep.index + 1) / VerificationStep.values.length * 100).round();
  }

  double _getProgressValue() {
    return (currentStep.index + 1) / VerificationStep.values.length;
  }

  String _getStepTitle() {
    switch (currentStep) {
      case VerificationStep.personalInfo:
        return 'Personal Information';
      case VerificationStep.idUpload:
        return 'ID Verification';
      case VerificationStep.certifications:
        return 'Certifications';
      case VerificationStep.backgroundCheck:
        return 'Background Check';
      case VerificationStep.review:
        return 'Review & Submit';
    }
  }

  String _getStepDescription() {
    switch (currentStep) {
      case VerificationStep.personalInfo:
        return 'Please provide your basic information';
      case VerificationStep.idUpload:
        return 'Upload your ID documents for verification';
      case VerificationStep.certifications:
        return 'Select your professional certifications';
      case VerificationStep.backgroundCheck:
        return 'Confirm your background and agree to terms';
      case VerificationStep.review:
        return 'Review all information before submitting';
    }
  }

  void _showSnackBar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: const Color(0xFF34C759),
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
      ),
    );
  }
} 