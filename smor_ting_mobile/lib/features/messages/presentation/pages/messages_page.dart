import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

class Conversation {
  final String id;
  final String providerName;
  final String providerImage;
  final String lastMessage;
  final DateTime timestamp;
  int unreadCount;
  final bool isOnline;

  Conversation({
    required this.id,
    required this.providerName,
    required this.providerImage,
    required this.lastMessage,
    required this.timestamp,
    this.unreadCount = 0,
    this.isOnline = false,
  });
}

class MessagesPage extends ConsumerStatefulWidget {
  const MessagesPage({super.key});

  @override
  ConsumerState<MessagesPage> createState() => _MessagesPageState();
}

class _MessagesPageState extends ConsumerState<MessagesPage> {
  List<Conversation> conversations = [
    Conversation(
      id: '1',
      providerName: 'John Smith',
      providerImage: '👨‍🔧',
      lastMessage: 'I\'ll be there in 15 minutes',
      timestamp: DateTime.now().subtract(const Duration(minutes: 5)),
      unreadCount: 2,
      isOnline: true,
    ),
    Conversation(
      id: '2',
      providerName: 'Sarah Wilson',
      providerImage: '👩‍🔧',
      lastMessage: 'The service has been completed',
      timestamp: DateTime.now().subtract(const Duration(hours: 2)),
      unreadCount: 0,
      isOnline: false,
    ),
    Conversation(
      id: '3',
      providerName: 'Mike Johnson',
      providerImage: '👨‍🔧',
      lastMessage: 'Can you confirm the address?',
      timestamp: DateTime.now().subtract(const Duration(days: 1)),
      unreadCount: 1,
      isOnline: true,
    ),
    Conversation(
      id: '4',
      providerName: 'David Brown',
      providerImage: '👨‍🔧',
      lastMessage: 'Thank you for choosing our service',
      timestamp: DateTime.now().subtract(const Duration(days: 2)),
      unreadCount: 0,
      isOnline: false,
    ),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF2F2F7),
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Color(0xFF007AFF)),
          onPressed: () => context.pop(),
        ),
        title: const Text(
          'Messages',
          style: TextStyle(
            color: Colors.black,
            fontSize: 18,
            fontWeight: FontWeight.w600,
          ),
        ),
        centerTitle: true,
        actions: [
          IconButton(
            icon: const Icon(Icons.search, color: Color(0xFF007AFF)),
            onPressed: () {
              _showSearchDialog();
            },
          ),
        ],
      ),
      body: conversations.isEmpty
          ? _buildEmptyState()
          : ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: conversations.length,
              itemBuilder: (context, index) {
                final conversation = conversations[index];
                return _buildConversationCard(conversation);
              },
            ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          _showNewMessageDialog();
        },
        backgroundColor: const Color(0xFFFF3B30),
        foregroundColor: Colors.white,
        child: const Icon(Icons.message),
      ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 80,
            height: 80,
            decoration: BoxDecoration(
              color: const Color(0xFFE5E5EA),
              borderRadius: BorderRadius.circular(40),
            ),
            child: const Icon(
              Icons.message,
              size: 40,
              color: Color(0xFF8E8E93),
            ),
          ),
          const SizedBox(height: 16),
          const Text(
            'No Messages',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.w600,
              color: Color(0xFF8E8E93),
            ),
          ),
          const SizedBox(height: 8),
          const Text(
                            'Start a conversation with your agents',
            style: TextStyle(
              fontSize: 14,
              color: Color(0xFF8E8E93),
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 24),
          ElevatedButton.icon(
            onPressed: () {
              _showNewMessageDialog();
            },
            icon: const Icon(Icons.message, size: 20),
            label: const Text('Start Conversation'),
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFFFF3B30),
              foregroundColor: Colors.white,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
              padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildConversationCard(Conversation conversation) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: const Color(0xFFE5E5EA)),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: ListTile(
        contentPadding: const EdgeInsets.all(16),
        leading: Stack(
          children: [
            Container(
              width: 56,
              height: 56,
              decoration: BoxDecoration(
                color: const Color(0xFFF2F2F7),
                borderRadius: BorderRadius.circular(28),
              ),
              child: Center(
                child: Text(
                  conversation.providerImage,
                  style: const TextStyle(fontSize: 24),
                ),
              ),
            ),
            if (conversation.isOnline)
              Positioned(
                bottom: 0,
                right: 0,
                child: Container(
                  width: 16,
                  height: 16,
                  decoration: BoxDecoration(
                    color: const Color(0xFF34C759),
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(color: Colors.white, width: 2),
                  ),
                ),
              ),
          ],
        ),
        title: Row(
          children: [
            Expanded(
              child: Text(
                conversation.providerName,
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                  color: Colors.black,
                ),
              ),
            ),
            Text(
              _formatTimestamp(conversation.timestamp),
              style: const TextStyle(
                fontSize: 12,
                color: Color(0xFF8E8E93),
              ),
            ),
          ],
        ),
        subtitle: Row(
          children: [
            Expanded(
              child: Text(
                conversation.lastMessage,
                style: TextStyle(
                  fontSize: 14,
                  color: conversation.unreadCount > 0
                      ? Colors.black
                      : const Color(0xFF8E8E93),
                  fontWeight: conversation.unreadCount > 0
                      ? FontWeight.w500
                      : FontWeight.normal,
                ),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ),
            if (conversation.unreadCount > 0) ...[
              const SizedBox(width: 8),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                decoration: BoxDecoration(
                  color: const Color(0xFFFF3B30),
                  borderRadius: BorderRadius.circular(10),
                ),
                child: Text(
                  conversation.unreadCount.toString(),
                  style: const TextStyle(
                    fontSize: 10,
                    color: Colors.white,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
            ],
          ],
        ),
        onTap: () {
          _openChat(conversation);
        },
        trailing: PopupMenuButton<String>(
          icon: const Icon(Icons.more_vert, color: Color(0xFF8E8E93)),
          onSelected: (value) {
            _handleConversationAction(value, conversation);
          },
          itemBuilder: (context) => [
            const PopupMenuItem(
              value: 'mark_read',
              child: Row(
                children: [
                  Icon(Icons.mark_email_read, color: Color(0xFF007AFF)),
                  SizedBox(width: 8),
                  Text('Mark as Read'),
                ],
              ),
            ),
            const PopupMenuItem(
              value: 'delete',
              child: Row(
                children: [
                  Icon(Icons.delete, color: Color(0xFFFF3B30)),
                  SizedBox(width: 8),
                  Text('Delete'),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  String _formatTimestamp(DateTime timestamp) {
    final now = DateTime.now();
    final difference = now.difference(timestamp);

    if (difference.inMinutes < 1) {
      return 'Now';
    } else if (difference.inMinutes < 60) {
      return '${difference.inMinutes}m';
    } else if (difference.inHours < 24) {
      return '${difference.inHours}h';
    } else if (difference.inDays < 7) {
      return '${difference.inDays}d';
    } else {
      return '${timestamp.day}/${timestamp.month}';
    }
  }

  void _openChat(Conversation conversation) {
    // Navigate to chat screen
    _showChatDialog(conversation);
  }

  void _handleConversationAction(String action, Conversation conversation) {
    switch (action) {
      case 'mark_read':
        setState(() {
          conversation.unreadCount = 0;
        });
        _showSnackBar('Marked as read');
        break;
      case 'delete':
        _showDeleteConfirmationDialog(conversation);
        break;
    }
  }

  void _showSearchDialog() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Search Messages'),
          content: const Text('Search functionality will be implemented here.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.of(context).pop();
                _showSnackBar('Search feature coming soon');
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFFF3B30),
                foregroundColor: Colors.white,
              ),
              child: const Text('Search'),
            ),
          ],
        );
      },
    );
  }

  void _showNewMessageDialog() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('New Message'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Text('Select a provider to start a conversation:'),
              const SizedBox(height: 16),
              _buildProviderOption('John Smith', '👨‍🔧', () {
                Navigator.of(context).pop();
                _startNewConversation('John Smith');
              }),
              const SizedBox(height: 8),
              _buildProviderOption('Sarah Wilson', '👩‍🔧', () {
                Navigator.of(context).pop();
                _startNewConversation('Sarah Wilson');
              }),
              const SizedBox(height: 8),
              _buildProviderOption('Mike Johnson', '👨‍🔧', () {
                Navigator.of(context).pop();
                _startNewConversation('Mike Johnson');
              }),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
          ],
        );
      },
    );
  }

  Widget _buildProviderOption(String name, String emoji, VoidCallback onTap) {
    return InkWell(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          border: Border.all(color: const Color(0xFFE5E5EA)),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Row(
          children: [
            Text(emoji, style: const TextStyle(fontSize: 20)),
            const SizedBox(width: 12),
            Text(
              name,
              style: const TextStyle(
                fontSize: 16,
                color: Colors.black,
              ),
            ),
            const Spacer(),
            const Icon(Icons.arrow_forward_ios, color: Color(0xFF8E8E93), size: 16),
          ],
        ),
      ),
    );
  }

  void _startNewConversation(String providerName) {
    // Add new conversation to the list
    final newConversation = Conversation(
      id: DateTime.now().millisecondsSinceEpoch.toString(),
      providerName: providerName,
      providerImage: '👨‍🔧',
      lastMessage: 'New conversation started',
      timestamp: DateTime.now(),
      unreadCount: 0,
      isOnline: true,
    );

    setState(() {
      conversations.insert(0, newConversation);
    });

    _showSnackBar('Conversation started with $providerName');
  }

  void _showChatDialog(Conversation conversation) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: Text('Chat with ${conversation.providerName}'),
          content: const Text('Chat functionality will be implemented here.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Close'),
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.of(context).pop();
                _showSnackBar('Chat opened');
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFFF3B30),
                foregroundColor: Colors.white,
              ),
              child: const Text('Open Chat'),
            ),
          ],
        );
      },
    );
  }

  void _showDeleteConfirmationDialog(Conversation conversation) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Delete Conversation'),
          content: Text('Are you sure you want to delete the conversation with ${conversation.providerName}?'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                setState(() {
                  conversations.removeWhere((c) => c.id == conversation.id);
                });
                Navigator.of(context).pop();
                _showSnackBar('Conversation deleted');
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFFF3B30),
                foregroundColor: Colors.white,
              ),
              child: const Text('Delete'),
            ),
          ],
        );
      },
    );
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