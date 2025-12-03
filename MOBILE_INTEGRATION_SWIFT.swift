//
//  UnifiedChatScreen.swift
//  Lomi Dating App
//
//  UNIFIED CHAT SCREEN - Handles both private 1-on-1 chat AND TikTok-style live streaming chat
//  Same UI, same WebSocket connection, different modes
//

import SwiftUI
import Starscream

// MARK: - Chat Mode
enum ChatMode: String, Codable {
    case `private` = "private"
    case live = "live"
}

// MARK: - Message Model
struct ChatMessage: Identifiable, Codable {
    let id: String
    let type: String // "message", "join", "leave", "gift", "pin", "system"
    let mode: ChatMode
    let content: String?
    let messageType: String? // "text", "photo", "video", "gift", "system"
    let senderID: String?
    let senderName: String?
    let senderAvatar: String?
    let receiverID: String?
    let timestamp: String
    
    // Live chat specific
    let liveStreamID: String?
    let seq: Int64?
    let viewerCount: Int?
    let isPinned: Bool?
    let isSystem: Bool?
    
    // Private chat specific
    let matchID: String?
    
    // Common
    let mediaURL: String?
    let giftID: String?
    let metadata: [String: AnyCodable]?
    
    enum CodingKeys: String, CodingKey {
        case id = "message_id"
        case type, mode, content
        case messageType = "message_type"
        case senderID = "sender_id"
        case senderName = "sender_name"
        case senderAvatar = "sender_avatar"
        case receiverID = "receiver_id"
        case timestamp
        case liveStreamID = "live_stream_id"
        case seq
        case viewerCount = "viewer_count"
        case isPinned = "is_pinned"
        case isSystem = "is_system"
        case matchID = "match_id"
        case mediaURL = "media_url"
        case giftID = "gift_id"
        case metadata
    }
}

// MARK: - Unified Chat Screen
struct UnifiedChatScreen: View {
    // MARK: - Properties
    let mode: ChatMode
    let matchID: String? // For private chat
    let liveStreamID: String? // For live chat
    let isBroadcaster: Bool // True if user owns the live stream
    
    @StateObject private var viewModel: ChatViewModel
    @State private var messageText = ""
    @State private var showGiftPicker = false
    @State private var pinnedMessage: ChatMessage?
    
    // MARK: - Initialization
    init(mode: ChatMode, matchID: String? = nil, liveStreamID: String? = nil, isBroadcaster: Bool = false) {
        self.mode = mode
        self.matchID = matchID
        self.liveStreamID = liveStreamID
        self.isBroadcaster = isBroadcaster
        
        _viewModel = StateObject(wrappedValue: ChatViewModel(
            mode: mode,
            matchID: matchID,
            liveStreamID: liveStreamID,
            isBroadcaster: isBroadcaster
        ))
    }
    
    // MARK: - Body
    var body: some View {
        VStack(spacing: 0) {
            // Header
            headerView
            
            // Pinned message (live mode only)
            if mode == .live, let pinned = pinnedMessage {
                pinnedMessageView(pinned)
            }
            
            // Messages
            ScrollViewReader { proxy in
                ScrollView {
                    LazyVStack(spacing: 8) {
                        ForEach(viewModel.messages) { message in
                            messageRow(message)
                        }
                    }
                    .padding()
                }
                .onChange(of: viewModel.messages.count) { _ in
                    if let lastMessage = viewModel.messages.last {
                        withAnimation {
                            proxy.scrollTo(lastMessage.id, anchor: .bottom)
                        }
                    }
                }
            }
            
            // Input bar
            inputBar
        }
        .onAppear {
            viewModel.connect()
        }
        .onDisappear {
            viewModel.disconnect()
        }
        .sheet(isPresented: $showGiftPicker) {
            GiftPickerView { gift in
                viewModel.sendGift(gift)
                showGiftPicker = false
            }
        }
    }
    
    // MARK: - Header View
    private var headerView: some View {
        HStack {
            if mode == .private {
                // Private chat header
                Text("Chat")
                    .font(.headline)
            } else {
                // Live chat header
                HStack(spacing: 8) {
                    Image(systemName: "eye.fill")
                        .foregroundColor(.red)
                    Text("\(viewModel.viewerCount)")
                        .font(.subheadline)
                        .fontWeight(.semibold)
                }
            }
            
            Spacer()
            
            Button(action: {
                showGiftPicker = true
            }) {
                Image(systemName: "gift.fill")
                    .foregroundColor(.pink)
            }
        }
        .padding()
        .background(Color(.systemBackground))
        .shadow(radius: 2)
    }
    
    // MARK: - Pinned Message View
    private func pinnedMessageView(_ message: ChatMessage) -> some View {
        HStack {
            Image(systemName: "pin.fill")
                .foregroundColor(.orange)
            
            Text(message.content ?? "")
                .font(.caption)
                .lineLimit(1)
            
            Spacer()
            
            if isBroadcaster {
                Button(action: {
                    viewModel.unpinMessage()
                }) {
                    Image(systemName: "xmark")
                        .foregroundColor(.gray)
                }
            }
        }
        .padding(.horizontal)
        .padding(.vertical, 8)
        .background(Color.orange.opacity(0.1))
    }
    
    // MARK: - Message Row
    private func messageRow(_ message: ChatMessage) -> some View {
        Group {
            if message.type == "join" || message.type == "leave" || message.isSystem == true {
                // System message (centered)
                systemMessageView(message)
            } else if message.type == "gift" {
                // Gift message with animation
                giftMessageView(message)
            } else {
                // Regular message
                regularMessageView(message)
            }
        }
    }
    
    private func systemMessageView(_ message: ChatMessage) -> some View {
        Text(message.content ?? "")
            .font(.caption)
            .foregroundColor(.gray)
            .padding(.horizontal, 12)
            .padding(.vertical, 6)
            .background(Color.gray.opacity(0.1))
            .cornerRadius(12)
            .frame(maxWidth: .infinity)
    }
    
    private func giftMessageView(_ message: ChatMessage) -> some View {
        HStack {
            if mode == .live {
                // Live mode: floating bubble style
                VStack(alignment: .leading, spacing: 4) {
                    Text(message.senderName ?? "Anonymous")
                        .font(.caption)
                        .fontWeight(.semibold)
                    
                    HStack {
                        Text("🎁")
                        Text("sent a gift!")
                            .font(.caption)
                    }
                }
                .padding(.horizontal, 12)
                .padding(.vertical, 8)
                .background(Color.pink.opacity(0.2))
                .cornerRadius(16)
            } else {
                // Private mode: regular message style
                VStack(alignment: .leading, spacing: 4) {
                    Text(message.senderName ?? "")
                        .font(.caption)
                        .fontWeight(.semibold)
                    
                    HStack {
                        Text("🎁 Gift")
                        Spacer()
                    }
                    .padding()
                    .background(Color.pink.opacity(0.1))
                    .cornerRadius(12)
                }
            }
            
            Spacer()
        }
    }
    
    private func regularMessageView(_ message: ChatMessage) -> some View {
        HStack {
            if mode == .live {
                // Live mode: floating bubble with avatar
                HStack(alignment: .top, spacing: 8) {
                    // Avatar
                    AsyncImage(url: URL(string: message.senderAvatar ?? "")) { image in
                        image.resizable()
                    } placeholder: {
                        Circle()
                            .fill(Color.gray.opacity(0.3))
                    }
                    .frame(width: 32, height: 32)
                    .clipShape(Circle())
                    
                    // Message bubble
                    VStack(alignment: .leading, spacing: 2) {
                        Text(message.senderName ?? "Anonymous")
                            .font(.caption)
                            .fontWeight(.semibold)
                            .foregroundColor(.orange)
                        
                        Text(message.content ?? "")
                            .font(.subheadline)
                    }
                    .padding(.horizontal, 12)
                    .padding(.vertical, 8)
                    .background(Color.black.opacity(0.6))
                    .foregroundColor(.white)
                    .cornerRadius(16)
                }
                
                Spacer()
                
                // Pin button (broadcaster only)
                if isBroadcaster {
                    Button(action: {
                        viewModel.pinMessage(message)
                    }) {
                        Image(systemName: "pin")
                            .foregroundColor(.gray)
                    }
                }
            } else {
                // Private mode: traditional chat bubbles
                let isCurrentUser = message.senderID == viewModel.currentUserID
                
                if isCurrentUser {
                    Spacer()
                }
                
                VStack(alignment: isCurrentUser ? .trailing : .leading, spacing: 4) {
                    if !isCurrentUser {
                        Text(message.senderName ?? "")
                            .font(.caption)
                            .foregroundColor(.gray)
                    }
                    
                    Text(message.content ?? "")
                        .padding(.horizontal, 12)
                        .padding(.vertical, 8)
                        .background(isCurrentUser ? Color.blue : Color.gray.opacity(0.2))
                        .foregroundColor(isCurrentUser ? .white : .primary)
                        .cornerRadius(16)
                }
                
                if !isCurrentUser {
                    Spacer()
                }
            }
        }
    }
    
    // MARK: - Input Bar
    private var inputBar: some View {
        HStack(spacing: 12) {
            TextField(mode == .live ? "Say something..." : "Message...", text: $messageText)
                .textFieldStyle(RoundedBorderTextFieldStyle())
                .onChange(of: messageText) { _ in
                    if mode == .private {
                        viewModel.sendTypingIndicator()
                    }
                }
            
            Button(action: sendMessage) {
                Image(systemName: "paperplane.fill")
                    .foregroundColor(messageText.isEmpty ? .gray : .blue)
            }
            .disabled(messageText.isEmpty)
        }
        .padding()
        .background(Color(.systemBackground))
    }
    
    // MARK: - Actions
    private func sendMessage() {
        guard !messageText.isEmpty else { return }
        viewModel.sendMessage(messageText)
        messageText = ""
    }
}

// MARK: - Chat View Model
class ChatViewModel: ObservableObject {
    @Published var messages: [ChatMessage] = []
    @Published var viewerCount: Int = 0
    @Published var isConnected: Bool = false
    
    let mode: ChatMode
    let matchID: String?
    let liveStreamID: String?
    let isBroadcaster: Bool
    let currentUserID: String
    
    private var socket: WebSocket?
    private var lastSeq: Int64 = 0
    
    init(mode: ChatMode, matchID: String?, liveStreamID: String?, isBroadcaster: Bool) {
        self.mode = mode
        self.matchID = matchID
        self.liveStreamID = liveStreamID
        self.isBroadcaster = isBroadcaster
        self.currentUserID = UserDefaults.standard.string(forKey: "user_id") ?? ""
    }
    
    // MARK: - WebSocket Connection
    func connect() {
        guard let token = UserDefaults.standard.string(forKey: "auth_token") else { return }
        
        var urlString = "wss://api.lomi.app/ws/chat?token=\(token)&mode=\(mode.rawValue)"
        
        if mode == .private, let matchID = matchID {
            urlString += "&match_id=\(matchID)"
        } else if mode == .live, let liveStreamID = liveStreamID {
            urlString += "&live_stream_id=\(liveStreamID)"
            urlString += "&is_broadcaster=\(isBroadcaster)"
            if lastSeq > 0 {
                urlString += "&last_seq=\(lastSeq)"
            }
        }
        
        guard let url = URL(string: urlString) else { return }
        
        var request = URLRequest(url: url)
        request.timeoutInterval = 5
        
        socket = WebSocket(request: request)
        socket?.delegate = self
        socket?.connect()
    }
    
    func disconnect() {
        socket?.disconnect()
        socket = nil
    }
    
    // MARK: - Send Messages
    func sendMessage(_ text: String) {
        var message: [String: Any] = [
            "type": "message",
            "mode": mode.rawValue,
            "message_type": "text",
            "content": text,
            "timestamp": ISO8601DateFormatter().string(from: Date())
        ]
        
        if mode == .private, let matchID = matchID {
            message["match_id"] = matchID
        } else if mode == .live, let liveStreamID = liveStreamID {
            message["live_stream_id"] = liveStreamID
        }
        
        sendJSON(message)
    }
    
    func sendGift(_ gift: Gift) {
        var message: [String: Any] = [
            "type": "gift",
            "mode": mode.rawValue,
            "message_type": "gift",
            "gift_id": gift.id,
            "timestamp": ISO8601DateFormatter().string(from: Date())
        ]
        
        if mode == .private, let matchID = matchID {
            message["match_id"] = matchID
        } else if mode == .live, let liveStreamID = liveStreamID {
            message["live_stream_id"] = liveStreamID
        }
        
        sendJSON(message)
    }
    
    func sendTypingIndicator() {
        guard mode == .private, let matchID = matchID else { return }
        
        let message: [String: Any] = [
            "type": "typing",
            "mode": "private",
            "match_id": matchID,
            "is_typing": true,
            "timestamp": ISO8601DateFormatter().string(from: Date())
        ]
        
        sendJSON(message)
    }
    
    func pinMessage(_ message: ChatMessage) {
        guard mode == .live, isBroadcaster, let liveStreamID = liveStreamID else { return }
        
        let pinMessage: [String: Any] = [
            "type": "pin",
            "mode": "live",
            "live_stream_id": liveStreamID,
            "message_id": message.id,
            "content": message.content ?? "",
            "timestamp": ISO8601DateFormatter().string(from: Date())
        ]
        
        sendJSON(pinMessage)
    }
    
    func unpinMessage() {
        // Implement unpin logic
    }
    
    private func sendJSON(_ data: [String: Any]) {
        guard let jsonData = try? JSONSerialization.data(withJSONObject: data),
              let jsonString = String(data: jsonData, encoding: .utf8) else { return }
        socket?.write(string: jsonString)
    }
}

// MARK: - WebSocket Delegate
extension ChatViewModel: WebSocketDelegate {
    func didReceive(event: WebSocketEvent, client: WebSocket) {
        switch event {
        case .connected(_):
            DispatchQueue.main.async {
                self.isConnected = true
            }
            
        case .disconnected(_, _):
            DispatchQueue.main.async {
                self.isConnected = false
            }
            
        case .text(let string):
            handleMessage(string)
            
        case .binary(_):
            break
            
        case .error(let error):
            print("WebSocket error: \(error?.localizedDescription ?? "unknown")")
            
        default:
            break
        }
    }
    
    private func handleMessage(_ text: String) {
        guard let data = text.data(using: .utf8),
              let message = try? JSONDecoder().decode(ChatMessage.self, from: data) else { return }
        
        DispatchQueue.main.async {
            // Update viewer count
            if let count = message.viewerCount {
                self.viewerCount = count
            }
            
            // Update last sequence
            if let seq = message.seq {
                self.lastSeq = max(self.lastSeq, seq)
            }
            
            // Add message to list
            if message.type == "message" || message.type == "gift" || message.type == "join" || message.type == "leave" {
                self.messages.append(message)
            }
        }
    }
}

// MARK: - Helper Types
struct Gift: Identifiable {
    let id: String
    let name: String
    let coinPrice: Int
}

struct GiftPickerView: View {
    let onSelect: (Gift) -> Void
    
    var body: some View {
        // Implement gift picker UI
        Text("Gift Picker")
    }
}

struct AnyCodable: Codable {
    let value: Any
    
    init(_ value: Any) {
        self.value = value
    }
    
    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        if let string = try? container.decode(String.self) {
            value = string
        } else if let int = try? container.decode(Int.self) {
            value = int
        } else if let double = try? container.decode(Double.self) {
            value = double
        } else if let bool = try? container.decode(Bool.self) {
            value = bool
        } else {
            value = ""
        }
    }
    
    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        if let string = value as? String {
            try container.encode(string)
        } else if let int = value as? Int {
            try container.encode(int)
        } else if let double = value as? Double {
            try container.encode(double)
        } else if let bool = value as? Bool {
            try container.encode(bool)
        }
    }
}

// MARK: - Usage Examples

// Example 1: Private Chat
// let chatScreen = UnifiedChatScreen(
//     mode: .private,
//     matchID: "123e4567-e89b-12d3-a456-426614174000"
// )

// Example 2: Live Chat (Viewer)
// let chatScreen = UnifiedChatScreen(
//     mode: .live,
//     liveStreamID: "789e4567-e89b-12d3-a456-426614174000",
//     isBroadcaster: false
// )

// Example 3: Live Chat (Broadcaster)
// let chatScreen = UnifiedChatScreen(
//     mode: .live,
//     liveStreamID: "789e4567-e89b-12d3-a456-426614174000",
//     isBroadcaster: true
// )
